package main

import (
	"log"
	"os"

	"../database"
	"../helpers"
	"../models/boardPack"
	"../models/cardPack"
	"../models/devicePack"
	"../models/eventPack"
	"../ws"
)

func main() {
	env := os.Getenv("ENV")
	database.Initialize(env)

	boards, errGetBoards := boardPack.GetAllByDeadline()
	if errGetBoards != nil {
		helpers.DaemonLogToFile(helpers.Join(errGetBoards.Error(), " daemon:b_notifier"))
		return
	}

	cards, errGetCards := cardPack.GetAllByDeadline()
	if errGetCards != nil {
		helpers.DaemonLogToFile(helpers.Join(errGetCards.Error(), " daemon:b_notifier"))
		return
	}

	_, erTx := database.BeginTx()
	if erTx != nil {
		helpers.DaemonLogToFile(helpers.Join(erTx.Error(), " daemon:b_notifier"))
		database.GetConnection().Tx.Rollback()
		return
	}

	prepareBoards(boards)
	prepareCards(cards)

	database.GetConnection().Tx.Commit()
}

func prepareBoards(boards boardPack.Boards) {
	for _, v := range boards {
		d, e := devicePack.GetByUserId(v.Pm)
		if e != nil {
			log.Println(e.Error())
			helpers.DaemonLogToFile(helpers.Join(e.Error(), " daemon:b_notifier"))
			database.GetConnection().Tx.Rollback()
			break
		}

		event := eventPack.Event{
			EventName: ws.EVENTS_MAP["BOARD_DEADLINE"],
			Data:      v.Render("leadership"),
			DeviceId:  d.Id,
		}
		event.SetDefault()
		event.SetEventType(eventPack.TYPE[2])
		err := event.Create()
		if err != nil {
			log.Println(err.Error())
			database.GetConnection().Tx.Rollback()
			break
		}

		v.SetDaemonStatus("parse")
		_, errBoardUpdate := boardPack.Update(&v)
		if errBoardUpdate != nil {
			log.Println(errBoardUpdate.Error())
			helpers.DaemonLogToFile(helpers.Join(errBoardUpdate.Error(), " daemon:b_notifier"))
			database.GetConnection().Tx.Rollback()
			break
		}
	}
}

func prepareCards(cards cardPack.Cards) {
	for _, v := range cards {

		log.Println(v.IdMembers)
		if len(v.IdMembers) > 0 {
			devices, e := devicePack.GetByUserIds(v.IdMembers)
			if e != nil {
				log.Println(e.Error())
				helpers.DaemonLogToFile(helpers.Join(e.Error(), " daemon:b_notifier"))
				database.GetConnection().Tx.Rollback()
				break
			}

			if len(devices) > 0 {
				for _, d := range devices {
					event := eventPack.Event{
						EventName: ws.EVENTS_MAP["CARD_DEADLINE"],
						Data:      v.Render("leadership"),
						DeviceId:  d.Id,
					}
					event.SetDefault()
					event.SetEventType(eventPack.TYPE[2])
					err := event.Create()
					if err != nil {
						log.Println(err.Error())
						database.GetConnection().Tx.Rollback()
						break
					}

					v.SetDaemonStatus("parse")
					_, errBoardUpdate := cardPack.Update(&v)
					if errBoardUpdate != nil {
						log.Println(errBoardUpdate.Error())
						helpers.DaemonLogToFile(helpers.Join(errBoardUpdate.Error(), " daemon:b_notifier"))
						database.GetConnection().Tx.Rollback()
						break
					}
				}
			}
		}
	}
}
