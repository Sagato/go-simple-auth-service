package api

func (s *server) routes() {
	s.router.HandleFunc("/register", s.RegisterUser).Methods("POST")
}
