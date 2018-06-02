package p2p

import (
	"errors"

	"github.com/gizo-network/gizo/helpers"

	"github.com/dghubble/sling"
	"github.com/kpango/glg"
)

var (
	s = sling.New().Base(CentrumURL).Add("User-Agent", "Gizo Node")
	//ErrNoToken occurs when node has not token
	ErrNoToken = errors.New("Centrum: No token in struct")
)

type (
	//DispatcherBody message sent to centrum
	DispatcherBody struct {
		Pub  string `url:"pub"`
		IP   string `url:"ip"`
		Port int    `url:"port"`
	}

	//Centrum communication with centrum
	Centrum struct {
		token  string // token received from centrum
		logger *glg.Glg
	}
)

func NewCentrum() *Centrum {
	return &Centrum{logger: helpers.Logger()}
}

//GetToken returns token
func (c Centrum) GetToken() string {
	return c.token
}

//SetToken sets dispatcher token
func (c *Centrum) SetToken(token string) {
	c.token = token
}

//GetDispatchers returns active dispatchers
func (c Centrum) GetDispatchers() (map[string]interface{}, error) {
	var dispatchers []string
	temp := make(map[string]interface{})
	_, err := s.New().Get("/v1/dispatchers").Receive(&dispatchers, &temp)
	if err != nil {
		return nil, err
	}
	if len(dispatchers) != 0 {
		temp["dispatchers"] = dispatchers
	}
	return temp, nil
}

//NewDisptcher registers dispatcher in centrum
func (c *Centrum) NewDisptcher(pub, ip string, port int) error {
	data := DispatcherBody{Pub: pub, IP: ip, Port: port}
	res := make(map[string]interface{})
	_, err := s.Post("/v1/dispatcher").BodyForm(data).Receive(&res, &res)
	if err != nil {
		return err
	}
	token, ok := res["token"]
	if !ok {
		return errors.New(res["status"].(string))
	}
	c.SetToken(token.(string))
	return nil
}

//ConnectWorker increments dispatchers worker in centrum
func (c Centrum) ConnectWorker() (map[string]interface{}, error) {
	if c.GetToken() == "" {
		return nil, ErrNoToken
	}
	res := make(map[string]interface{})
	_, err := s.Patch("/v1/dispatcher/connect").Set("x-gizo-token", c.GetToken()).Receive(&res, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//DisconnectWorker decrements dispatchers worker in centrum
func (c Centrum) DisconnectWorker() (map[string]interface{}, error) {
	if c.GetToken() == "" {
		return nil, ErrNoToken
	}
	res := make(map[string]interface{})
	_, err := s.Patch("/v1/dispatcher/disconnect").Set("x-gizo-token", c.GetToken()).Receive(&res, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//Wake changes node status to active in centrum
func (c Centrum) Wake() (map[string]interface{}, error) {
	//TODO: create new worker if wake fails
	if c.GetToken() == "" {
		return nil, ErrNoToken
	}
	res := make(map[string]interface{})
	_, err := s.Patch("/v1/dispatcher/wake").Set("x-gizo-token", c.GetToken()).Receive(&res, &res)
	if err != nil {
		return nil, err
	}
	c.logger.Log("Centrum: waking node")
	return res, nil
}

//Sleep changes node status to sleep in centrum
func (c Centrum) Sleep() (map[string]interface{}, error) {
	if c.GetToken() == "" {
		return nil, ErrNoToken
	}
	res := make(map[string]interface{})
	_, err := s.Patch("/v1/dispatcher/sleep").Set("x-gizo-token", c.GetToken()).Receive(&res, &res)
	if err != nil {
		return nil, err
	}
	c.logger.Log("Centrum: sleeping node")
	return res, nil
}
