package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/polyk005/tg_bot/internal/domain"
	"github.com/polyk005/tg_bot/internal/service"
	"github.com/polyk005/tg_bot/pkg/logger"
)

// UserState хранит состояние пользователя при вводе расписания
type UserState struct {
	CurrentDay     time.Weekday
	CurrentStep    int // 0 - ожидание дня, 1 - ожидание предметов, 2 - выбор недели
	ScheduleInput  map[time.Weekday][]domain.Lesson
	CurrentWeekNum int // 1 - числитель, 2 - знаменатель
}

var userStates = make(map[int64]*UserState)

// Добавляем в начало файла
func getCurrentWeekNumber() int {
	_, week := time.Now().ISOWeek()
	if week%2 == 0 {
		return 2 // Знаменатель
	}
	return 1 // Числитель
}

// RegisterHandlers регистрирует все обработчики команд бота
func RegisterHandlers(b *bot.Bot, svc *service.Service, log logger.Logger) {
	// Основные команды
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler(svc, log))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler(log))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/setschedule", bot.MatchTypeExact, setScheduleHandler(svc, log))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/schedule", bot.MatchTypeExact, scheduleHandler(svc, log))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/ask", bot.MatchTypePrefix, askHandler(svc, log))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cancel", bot.MatchTypeExact, cancelHandler(log))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/today", bot.MatchTypeExact, todayHandler(svc, log))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/tomorrow", bot.MatchTypeExact, tomorrowHandler(svc, log))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/week", bot.MatchTypeExact, weekHandler(svc, log))

	// Обработчик для любых текстовых сообщений (для пошагового ввода)
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, textMessageHandler(svc, log))
}

// startHandler обрабатывает команду /start
func startHandler(svc *service.Service, log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		if err := svc.ProcessStartCommand(ctx, userID); err != nil {
			log.Errorw("Failed to process start command", "error", err, "userID", userID)
			sendErrorMessage(b, ctx, chatID, "Не удалось начать работу. Попробуйте позже.")
			return
		}

		currentWeek := getCurrentWeekNumber()
		weekName := "числитель"
		if currentWeek == 2 {
			weekName = "знаменатель"
		}

		msg := fmt.Sprintf(`Привет! Я бот для управления расписанием. 
Сейчас %d-я учебная неделя (%s).

Основные команды:
/setschedule - установить расписание
/today - расписание на сегодня
/tomorrow - расписание на завтра
/week - расписание на неделю
/ask [вопрос] - задать вопрос AI-помощнику`, currentWeek, weekName)

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   msg,
		}); err != nil {
			log.Errorw("Failed to send welcome message", "error", err, "chatID", chatID)
		}
	}
}

// setScheduleHandler начинает процесс ввода расписания
func setScheduleHandler(svc *service.Service, log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		// Инициализируем состояние пользователя
		userStates[userID] = &UserState{
			ScheduleInput: make(map[time.Weekday][]domain.Lesson),
		}

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text: "📝 Введите расписание. Сначала укажите день недели (например, 'Понедельник')\n" +
				"Или отправьте /cancel для отмены",
			ReplyMarkup: weekdayKeyboard(),
		})
		if err != nil {
			log.Errorw("Failed to send schedule instructions", "error", err)
		}
	}
}

// cancelHandler отменяет ввод расписания
func cancelHandler(log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		delete(userStates, userID)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ Ввод расписания отменен",
		})
		if err != nil {
			log.Errorw("Failed to send cancel message", "error", err)
		}
	}
}

// textMessageHandler обрабатывает текстовые сообщения при вводе расписания
func textMessageHandler(svc *service.Service, log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID
		text := update.Message.Text

		state, exists := userStates[userID]
		if !exists {
			return // Не в режиме ввода расписания
		}

		switch state.CurrentStep {
		case 0: // Ожидаем день недели
			day, err := parseWeekday(text)
			if err != nil {
				_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:      chatID,
					Text:        "Неверный день недели. Пожалуйста, выберите день из списка:",
					ReplyMarkup: weekdayKeyboard(),
				})
				return
			}

			state.CurrentDay = day
			state.CurrentStep = 1

			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text: fmt.Sprintf("📅 Введите пары для %s в формате:\n"+
					"<Название> | <Начало> | <Конец> | <Аудитория> | <Преподаватель>\n\n"+
					"Пример:\n"+
					"Математика | 09:00 | 10:30 | Ауд. 101 | Иванов И.И.\n\n"+
					"Когда закончите, отправьте /done", russianWeekday(day)),
			})

		case 1: // Ожидаем пары
			if text == "/done" {
				// Завершаем ввод для этого дня
				_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text: fmt.Sprintf("День %s сохранен. Хотите добавить еще один день? (да/нет)",
						russianWeekday(state.CurrentDay)),
				})
				state.CurrentStep = 2 // Переходим к вопросу о продолжении
				return
			}

			lesson, err := parseLesson(text)
			if err != nil {
				_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "Неверный формат. Пожалуйста, используйте формат:\nНазвание | Начало | Конец | Аудитория | Преподаватель",
				})
				return
			}

			state.ScheduleInput[state.CurrentDay] = append(state.ScheduleInput[state.CurrentDay], lesson)
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "✅ Пара добавлена. Введите следующую или /done для завершения",
			})

		case 2: // Ожидаем ответ на вопрос о продолжении
			if strings.ToLower(text) == "да" {
				state.CurrentStep = 0
				_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:      chatID,
					Text:        "Выберите следующий день недели:",
					ReplyMarkup: weekdayKeyboard(),
				})
			} else {
				// Сохраняем расписание и завершаем
				for day, lessons := range state.ScheduleInput {
					svc.SaveSchedule(ctx, userID, day, lessons)
				}

				delete(userStates, userID)

				_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "✅ Расписание успешно сохранено!",
				})
			}
		}
	}
}

// todayHandler обрабатывает команду /today
func todayHandler(svc *service.Service, log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handleDaySchedule(ctx, b, update, svc, log, time.Now().Weekday())
	}
}

// tomorrowHandler обрабатывает команду /tomorrow
func tomorrowHandler(svc *service.Service, log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handleDaySchedule(ctx, b, update, svc, log, time.Now().Add(24*time.Hour).Weekday())
	}
}

// weekHandler обрабатывает команду /week
func weekHandler(svc *service.Service, log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID

		var sb strings.Builder
		sb.WriteString("📅 Ваше расписание на неделю:\n\n")

		for day := time.Monday; day <= time.Friday; day++ {
			lessons, err := svc.GetSchedule(ctx, userID, day)
			if err != nil {
				log.Errorw("Failed to get schedule", "error", err, "userID", userID, "day", day)
				continue
			}

			sb.WriteString(fmt.Sprintf("📌 %s:\n", russianWeekday(day)))

			if len(lessons) == 0 {
				sb.WriteString("    Пар нет\n\n")
				continue
			}

			for _, lesson := range lessons {
				sb.WriteString(fmt.Sprintf("    🕒 %s-%s %s\n",
					lesson.StartTime.Format("15:04"),
					lesson.EndTime.Format("15:04"),
					lesson.Name))
			}
			sb.WriteString("\n")
		}

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   sb.String(),
		}); err != nil {
			log.Errorw("Failed to send week schedule", "error", err, "chatID", chatID)
		}
	}
}

// handleDaySchedule обрабатывает расписание на конкретный день
func handleDaySchedule(ctx context.Context, b *bot.Bot, update *models.Update, svc *service.Service, log logger.Logger, day time.Weekday) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	lessons, err := svc.GetSchedule(ctx, userID, day)
	if err != nil {
		log.Errorw("Failed to get schedule", "error", err, "userID", userID, "day", day)
		sendErrorMessage(b, ctx, chatID, "Ошибка при получении расписания. Попробуйте позже.")
		return
	}

	if len(lessons) == 0 {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("На %s пар нет 🎉", russianWeekday(day)),
		}); err != nil {
			log.Errorw("Failed to send empty schedule message", "error", err, "chatID", chatID)
		}
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📅 Расписание на %s:\n\n", russianWeekday(day)))

	for i, lesson := range lessons {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, lesson.Name))
		sb.WriteString(fmt.Sprintf("   🕒 %s - %s\n", lesson.StartTime.Format("15:04"), lesson.EndTime.Format("15:04")))
		if lesson.Location != "" {
			sb.WriteString(fmt.Sprintf("   🏫 %s\n", lesson.Location))
		}
		if lesson.Teacher != "" {
			sb.WriteString(fmt.Sprintf("   👨‍🏫 %s\n", lesson.Teacher))
		}
		sb.WriteString("\n")
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   sb.String(),
	}); err != nil {
		log.Errorw("Failed to send schedule message", "error", err, "chatID", chatID)
	}
}

// helpHandler обрабатывает команду /help
func helpHandler(log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID

		msg := `ℹ️ Доступные команды:

/start - Начало работы
/help - Эта справка
/setschedule - Установить расписание
/today - Расписание на сегодня
/tomorrow - Расписание на завтра
/week - Расписание на неделю
/ask [вопрос] - Задать вопрос AI-помощнику`

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   msg,
		}); err != nil {
			log.Errorw("Failed to send help message", "error", err, "chatID", chatID)
		}
	}
}

// askHandler обрабатывает команду /ask
func askHandler(svc *service.Service, log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID
		question := strings.TrimPrefix(update.Message.Text, "/ask ")

		if question == "" {
			sendErrorMessage(b, ctx, chatID, "Пожалуйста, задайте вопрос после команды /ask")
			return
		}

		answer, err := svc.ProcessAIQuestion(ctx, question)
		if err != nil {
			log.Errorw("Failed to process AI question", "error", err, "question", question)
			sendErrorMessage(b, ctx, chatID, "Произошла ошибка при обработке вопроса.")
			return
		}

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   answer,
		}); err != nil {
			log.Errorw("Failed to send AI answer", "error", err, "chatID", chatID)
		}
	}
}

// scheduleHandler обрабатывает команду /schedule (показывает меню выбора)
func scheduleHandler(svc *service.Service, log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Сегодня", CallbackData: "schedule_today"},
					{Text: "Завтра", CallbackData: "schedule_tomorrow"},
				},
				{
					{Text: "Вся неделя", CallbackData: "schedule_week"},
				},
			},
		}

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatID,
			Text:        "Выберите период для просмотра расписания:",
			ReplyMarkup: kb,
		}); err != nil {
			log.Errorw("Failed to send schedule menu", "error", err, "chatID", chatID)
		}
	}
}

// Вспомогательные функции

// weekdayKeyboard создает клавиатуру с днями недели
func weekdayKeyboard() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "Понедельник"},
				{Text: "Вторник"},
				{Text: "Среда"},
			},
			{
				{Text: "Четверг"},
				{Text: "Пятница"},
				{Text: "Суббота"},
			},
			{
				{Text: "Воскресенье"},
				{Text: "/cancel"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
}

// parseWeekday преобразует текст в день недели
func parseWeekday(text string) (time.Weekday, error) {
	switch strings.ToLower(text) {
	case "понедельник":
		return time.Monday, nil
	case "вторник":
		return time.Tuesday, nil
	case "среда":
		return time.Wednesday, nil
	case "четверг":
		return time.Thursday, nil
	case "пятница":
		return time.Friday, nil
	case "суббота":
		return time.Saturday, nil
	case "воскресенье":
		return time.Sunday, nil
	default:
		return time.Sunday, fmt.Errorf("invalid weekday")
	}
}

// parseLesson разбирает строку с информацией о паре
func parseLesson(text string) (domain.Lesson, error) {
	parts := strings.Split(text, "|")
	if len(parts) != 5 {
		return domain.Lesson{}, fmt.Errorf("invalid format")
	}

	startTime, err := time.Parse("15:04", strings.TrimSpace(parts[1]))
	if err != nil {
		return domain.Lesson{}, err
	}

	endTime, err := time.Parse("15:04", strings.TrimSpace(parts[2]))
	if err != nil {
		return domain.Lesson{}, err
	}

	return domain.Lesson{
		Name:      strings.TrimSpace(parts[0]),
		StartTime: startTime,
		EndTime:   endTime,
		Location:  strings.TrimSpace(parts[3]),
		Teacher:   strings.TrimSpace(parts[4]),
	}, nil
}

// russianWeekday переводит день недели на русский
func russianWeekday(day time.Weekday) string {
	days := map[time.Weekday]string{
		time.Monday:    "Понедельник",
		time.Tuesday:   "Вторник",
		time.Wednesday: "Среда",
		time.Thursday:  "Четверг",
		time.Friday:    "Пятница",
		time.Saturday:  "Суббота",
		time.Sunday:    "Воскресенье",
	}
	return days[day]
}

// sendErrorMessage отправляет сообщение об ошибке
func sendErrorMessage(b *bot.Bot, ctx context.Context, chatID int64, text string) {
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "❌ " + text,
	}); err != nil {
		// Если не удалось отправить сообщение об ошибке,
		// здесь уже ничего не поделаешь, просто логируем
	}
}
