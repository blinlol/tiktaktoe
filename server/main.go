package main

import (
	"log"
	"net"

	"tiktaktoe/game"
	pb "tiktaktoe/game_proto"

	mapset "github.com/deckarep/golang-set/v2"
	"google.golang.org/grpc"
)


type server struct {
	pb.UnimplementedGameServer

	games []game.Game

	//ids of games that wait second player
	waiting_ids mapset.Set[int32]
}


func (s *server) findOrCreateGame() (new_game bool, game_id int32) {
	if s.waiting_ids.IsEmpty() {
		// new
		g := game.Game{
			Id: int32(len(s.games)),
			Field: game.GameField{},
			WaitChan: make(chan int),
			NextMovePlayer: pb.Player_CROSS,
			Streams: make(map[pb.Player]pb.Game_GameStreamServer)}
		s.games = append(s.games, g)
		new_game = true
		game_id = g.Id
		s.waiting_ids.Add(game_id)
	} else {
		new_game = false
		game_id, _ = s.waiting_ids.Pop()
	}
	return
}



func SendNextMove(game *game.Game, stream pb.Game_GameStreamServer) error {
	m := pb.GameServerMessage{
		Message: "",
		GameId: game.Id,
		Type: pb.ServerMessageType_MAKE_MOVE,
		Who: game.NextMovePlayer,
	}
	log.Println(&m)
	return stream.Send(&m)
}


func (s *server) GameStream(stream pb.Game_GameStreamServer) error {
	new_game, game_id := s.findOrCreateGame()
	log.Println(new_game, game_id)

	/*
	получили id игры в которую будет играть клиент
	если создана новая, то ждем подключения второго игрока
	если оба игрока найдены, то отправляем игрокам информацию об игре и том, кто они
	*/

	// ATOMIC !!!!!!!
	// одну игру обрабатывают 2 серверных процесса на каждый клиент
	game := &s.games[game_id]
	var iam pb.Player
	var err error

	if new_game {
		iam, err = game.InitFirstPlayer(stream)
		log.Println("init first")
	} else {
		iam, err = game.InitSecondPlayer(stream)
		log.Println("init second")
	}
	if err != nil {
		log.Fatalf("%v", err)
	}

	for {
		// сервер говорит всем, кто делает ход
		// err = game.SendNextMove()
		
		// ждет ответа
		if iam == game.NextMovePlayer {
			err = SendNextMove(game, stream)
			if err != nil {
				log.Fatalf("%v\n", err)
			}
			client_message, err := game.Streams[game.NextMovePlayer].Recv()
			log.Println(&client_message)
			if err != nil {
				log.Fatalf("%v\n", err)
			}
			// валидирует, что ход валидный
			row, col := client_message.Row, client_message.Col
			if ! game.Field.IsMoveValid(row, col) {
				continue
			}
			// сохраняет ход
			// отправляет ход всем
			end, winner, err := game.ApplyMove(row, col)
			if err != nil {
				log.Fatalf("%v\n", err)
			}
			if end {
				m := pb.GameServerMessage {
					Message: "End Game",
					GameId: game.Id,
					Type: pb.ServerMessageType_END_GAME,
					Who: winner,
				}
				game.Streams[pb.Player_CROSS].Send(&m)
				game.Streams[pb.Player_ZERO].Send(&m)
				break
			}
		}
		


		// клиенты получили свой паспорт
		// получили сообщение делать ход и если оно для них, то делают его
		// отправляют сделанный ход на сервер
		// получают от сервера сделанный ход и сохраняют его
	}
	return nil
}


func main(){
	lis, err := net.Listen("tcp", ":50052")
	if err != nil{
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	// pb.RegisterGameServer(s, &old_server{})
	pb.RegisterGameServer(s, &server{waiting_ids: mapset.NewSet[int32]()})

	log.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Serve error: %v\n", err)
	}
}