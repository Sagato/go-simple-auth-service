package api

func (s *server) routes() {
	s.router.HandleFunc("/register", s.RegisterUser).Methods("POST")
	s.router.HandleFunc("/activate", s.ActivateAccount).Methods("GET")
	s.router.HandleFunc("/login", s.Login).Methods("GET")
}
