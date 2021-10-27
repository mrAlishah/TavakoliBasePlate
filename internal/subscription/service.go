package subscription

func (s service) Test(req Request) (string, error){
	//implement login
	//msg := fmt.Sprintf("Hi %s", name)
	//if name != "Omid"{
	//	errMsg := errors.New("name is not equal to Omid")
	//	s.logger.Errorf("Test method err : %s", errMsg)
	//	return name, errMsg
	//}
	//return msg, nil
	return req.Email, nil
}
