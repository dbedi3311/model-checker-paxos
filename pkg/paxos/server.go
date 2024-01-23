package paxos

import (
	"coms4113/hw5/pkg/base"
)

const (
	Propose = "propose"
	Accept  = "accept"
	Decide  = "decide"
)

type Proposer struct {
	N             int
	Phase         string
	N_a_max       int
	V             interface{}
	SuccessCount  int
	ResponseCount int
	// To indicate if response from peer is received, should be initialized as []bool of len(server.peers)
	Responses []bool
	// Use this field to check if a message is latest.
	SessionId int

	// in case node will propose again - restore initial value
	InitialValue interface{}
}

type ServerAttribute struct {
	peers []base.Address
	me    int

	// Paxos parameter
	n_p int
	n_a int
	v_a interface{}

	// final result
	agreedValue interface{}

	// Propose parameter
	proposer Proposer

	// retry
	timeout *TimeoutTimer
}

type Server struct {
	base.CoreNode
	ServerAttribute
}

func NewServer(peers []base.Address, me int, proposedValue interface{}) *Server {
	response := make([]bool, len(peers))
	return &Server{
		CoreNode: base.CoreNode{},
		ServerAttribute: ServerAttribute{
			peers: peers,
			me:    me,
			proposer: Proposer{
				InitialValue: proposedValue,
				Responses:    response,
			},
			timeout: &TimeoutTimer{},
		},
	}
}

func (server *Server) MessageHandler(message base.Message) []base.Node {
	//TODO: implement it
	newNodes := make([]base.Node, 0, 1)
	newNode := server.copy()

	switch message.(type) {
	case *ProposeRequest:
		pr := message.(*ProposeRequest)
		n := pr.N

		sessionId := pr.SessionId

		ok := false
		if n > newNode.n_p {
			newNode.n_p = n
			ok = true
		}

		resp := &ProposeResponse{
			CoreMessage: base.MakeCoreMessage(message.To(), message.From()),
			Ok:          ok,
			SessionId:   sessionId,
			N_p:         newNode.n_p,
			N_a:         newNode.n_a,
			V_a:         newNode.v_a,
		}

		newNode.SetResponse([]base.Message{resp})
		newNodes = append(newNodes, newNode)

		return newNodes

	case *ProposeResponse:
		pr := message.(*ProposeResponse)

		ok := pr.Ok
		//n_p := pr.N_p
		n_a := pr.N_a
		v_a := pr.V_a

		sessionId := pr.SessionId

		// if newNode.proposer.V == "v1" && n_a == 1 {
		// 	fmt.Printf("v_a error: %+v\n", v_a)
		// 	panic("v1 n_a 1 is WRONG")
		// }

		if !base.IsNil(newNode.agreedValue) {
			newNodes = append(newNodes, newNode.copy())
			return newNodes
		}

		// if newNode.proposer.Phase != Propose {
		// 	newNodes = append(newNodes, newNode.copy())
		// 	return newNodes
		// }

		if sessionId != newNode.proposer.SessionId {
			newNodes = append(newNodes, newNode.copy())
			return newNodes
		}

		from := message.From()
		for i := range newNode.peers {
			if newNode.peers[i] == from {
				if newNode.proposer.Responses[i] == true {
					//newNode.proposer.ResponseCount -= 1
					newNodes = append(newNodes, newNode.copy())
					return newNodes
				}
				newNode.proposer.Responses[i] = true
			}
		}

		newNode.proposer.ResponseCount += 1

		if ok {
			// if newNode.proposer.V == "v1" {
			// 	fmt.Printf("server(%d): v1 is V   v_a: %+v  n_a: %d  N_a_max: %d  ok: %t, successcount: %d, from: %s\n", newNode.me, v_a, n_a, newNode.proposer.N_a_max, ok, newNode.proposer.SuccessCount, from)
			// }
			newNode.proposer.SuccessCount += 1

		}
		maj := len(newNode.peers)/2 + 1

		if n_a > newNode.proposer.N_a_max {
			// if ok {
			// 	fmt.Printf("proposer.V: %+v    proceed.v_a: %+v \n ", proceed.proposer.V, proceed.v_a)
			// 	panic("DEBUG ERROR REACHED")
			// }
			newNode.proposer.V = v_a
			newNode.proposer.N_a_max = n_a
		}

		if ok && newNode.proposer.SuccessCount >= maj {
			// two options: to proceed waiting for more responses, or move to next phase.
			proceed := newNode.copy()

			// if proceed.proposer.V == "v1" {
			// 	fmt.Printf("v1 is V   v_a: %+v  n_a: %d  N_a_max: %d  ok: %t, successcount: %d, from: %s\n", v_a, n_a, newNode.proposer.N_a_max, ok, newNode.proposer.SuccessCount, from)
			// 	panic("v1 accidentally set")
			// }

			proceed.proposer.Phase = Accept
			proceed.proposer.ResponseCount = 0
			proceed.proposer.SuccessCount = 0
			for i := range proceed.proposer.Responses {
				proceed.proposer.Responses[i] = false
			}

			response := make([]base.Message, 0, 1)

			for i := range proceed.peers {
				resp := &AcceptRequest{
					CoreMessage: base.MakeCoreMessage(proceed.Address(), proceed.peers[i]),
					N:           proceed.proposer.N,
					V:           proceed.proposer.V,
					SessionId:   sessionId,
				}
				response = append(response, resp)
			}

			proceed.SetResponse(response)
			newNodes = append(newNodes, proceed)

			// to wait for more responses
			newNodes = append(newNodes, newNode)

		} else {
			newNodes = append(newNodes, newNode)
		}

		return newNodes

	case *AcceptRequest:
		ar := message.(*AcceptRequest)
		n := ar.N
		val := ar.V
		sessionId := ar.SessionId

		// fmt.Printf("%+v \n", val)
		// fmt.Printf("%+v \n", n)
		// fmt.Printf("%+v \n", newNode.n_p)
		// panic("HELLO WORLD")

		// if val == "v1" && n > 0 {
		// 	panic("n_a being set improperly")
		// }

		ok := false
		if n >= newNode.n_p {
			newNode.n_p = n
			newNode.n_a = n
			newNode.v_a = val
			ok = true
		}

		resp := &AcceptResponse{
			CoreMessage: base.MakeCoreMessage(message.To(), message.From()),
			Ok:          ok,
			SessionId:   sessionId,
			N_p:         newNode.n_p,
		}

		newNode.SetResponse([]base.Message{resp})
		newNodes = append(newNodes, newNode)

		return newNodes

	case *AcceptResponse:
		ar := message.(*AcceptResponse)

		ok := ar.Ok
		//n_p := ar.N_p

		sessionId := ar.SessionId

		if !base.IsNil(newNode.agreedValue) {
			newNodes = append(newNodes, newNode.copy())
			return newNodes
		}

		if newNode.proposer.Phase != Accept {
			newNodes = append(newNodes, newNode.copy())
			return newNodes
		}

		if sessionId != newNode.proposer.SessionId {
			newNodes = append(newNodes, newNode)
			return newNodes
		}

		from := message.From()

		for i := range newNode.peers {
			if newNode.peers[i] == from {
				if newNode.proposer.Responses[i] == true {
					newNodes = append(newNodes, newNode)
					return newNodes
				}
				newNode.proposer.Responses[i] = true
			}
		}

		newNode.proposer.ResponseCount += 1

		maj := len(newNode.peers)/2 + 1

		if ok {
			newNode.proposer.SuccessCount += 1
		}

		if newNode.proposer.SuccessCount >= maj {
			// two options: to proceed waiting for more responses, or move to next phase.
			proceed := newNode.copy()

			proceed.proposer.Phase = Decide
			proceed.proposer.ResponseCount = 0
			proceed.proposer.SuccessCount = 0
			for i := range proceed.proposer.Responses {
				proceed.proposer.Responses[i] = false
			}

			response := make([]base.Message, 0, 1)

			for i := range proceed.peers {
				resp := &DecideRequest{
					CoreMessage: base.MakeCoreMessage(proceed.Address(), proceed.peers[i]),
					V:           proceed.proposer.V,
					SessionId:   sessionId,
				}
				response = append(response, resp)
			}

			proceed.SetResponse(response)
			newNodes = append(newNodes, proceed)

			// to wait for more responses
			newNodes = append(newNodes, newNode)

		} else {
			newNodes = append(newNodes, newNode)
		}

		return newNodes

	case *DecideRequest:
		dr := message.(*DecideRequest)

		val := dr.V
		sessionId := dr.SessionId

		ok := true
		newNode.v_a = val
		newNode.agreedValue = val
		newNode.proposer.InitialValue = val

		resp := &DecideResponse{
			CoreMessage: base.MakeCoreMessage(message.To(), message.From()),
			Ok:          ok,
			SessionId:   sessionId,
		}

		newNode.SetResponse([]base.Message{resp})
		newNodes = append(newNodes, newNode)

		return newNodes

	case *DecideResponse:
		dr := message.(*DecideResponse)

		//ok := dr.Ok
		sessionId := dr.SessionId

		if newNode.proposer.Phase != Decide {
			newNodes = append(newNodes, newNode.copy())
			return newNodes
		}

		if sessionId != newNode.proposer.SessionId {
			newNodes = append(newNodes, newNode.copy())
			return newNodes
		}

		newNodes = append(newNodes, newNode)
		return newNodes
	default:

		newNodes = append(newNodes, newNode)
		return newNodes
	}

	//panic("implement me")
}

// To start a new round of Paxos.
func (server *Server) StartPropose() {
	//TODO: implement it

	if base.IsNil(server.proposer.InitialValue) {
		return
	}

	newNode := server

	proposalNum := newNode.n_p + 1

	newNode.proposer.Phase = Propose
	newNode.proposer.ResponseCount = 0
	newNode.proposer.SuccessCount = 0
	newNode.proposer.Responses = []bool{false, false, false}
	newNode.proposer.N = proposalNum
	newNode.proposer.N_a_max = newNode.n_a
	newNode.proposer.SessionId += 1

	//fmt.Printf("server(%d): n_a - %d, \n", server.me, server.n_a)

	if base.IsNil(newNode.proposer.V) {
		newNode.proposer.V = newNode.proposer.InitialValue
	}

	response := make([]base.Message, 0, 1)

	for i := range newNode.peers {
		resp := &ProposeRequest{
			CoreMessage: base.MakeCoreMessage(newNode.Address(), newNode.peers[i]),
			N:           proposalNum,
			SessionId:   newNode.proposer.SessionId,
		}
		response = append(response, resp)
	}
	newNode.SetResponse(response)

	//newNode.n_p += 1

	//panic("implement me")
}

// Returns a deep copy of server node
func (server *Server) copy() *Server {
	response := make([]bool, len(server.peers))
	for i, flag := range server.proposer.Responses {
		response[i] = flag
	}

	var copyServer Server
	copyServer.me = server.me
	// shallow copy is enough, assuming it won't change
	copyServer.peers = server.peers
	copyServer.n_a = server.n_a
	copyServer.n_p = server.n_p
	copyServer.v_a = server.v_a
	copyServer.agreedValue = server.agreedValue
	copyServer.proposer = Proposer{
		N:             server.proposer.N,
		Phase:         server.proposer.Phase,
		N_a_max:       server.proposer.N_a_max,
		V:             server.proposer.V,
		SuccessCount:  server.proposer.SuccessCount,
		ResponseCount: server.proposer.ResponseCount,
		Responses:     response,
		InitialValue:  server.proposer.InitialValue,
		SessionId:     server.proposer.SessionId,
	}

	// doesn't matter, timeout timer is state-less
	copyServer.timeout = server.timeout

	return &copyServer
}

func (server *Server) NextTimer() base.Timer {
	return server.timeout
}

// A TimeoutTimer tick simulates the situation where a proposal procedure times out.
// It will close the current Paxos round and start a new one if no consensus reached so far,
// i.e. the server after timer tick will reset and restart from the first phase if Paxos not decided.
// The timer will not be activated if an agreed value is set.
func (server *Server) TriggerTimer() []base.Node {
	if server.timeout == nil {
		return nil
	}

	subNode := server.copy()
	subNode.StartPropose()

	return []base.Node{subNode}
}

func (server *Server) Attribute() interface{} {
	return server.ServerAttribute
}

func (server *Server) Copy() base.Node {
	return server.copy()
}

func (server *Server) Hash() uint64 {
	return base.Hash("paxos", server.ServerAttribute)
}

func (server *Server) Equals(other base.Node) bool {
	otherServer, ok := other.(*Server)

	if !ok || server.me != otherServer.me ||
		server.n_p != otherServer.n_p || server.n_a != otherServer.n_a || server.v_a != otherServer.v_a ||
		(server.timeout == nil) != (otherServer.timeout == nil) {
		return false
	}

	// fmt.Println(server.proposer.N != otherServer.proposer.N)
	// fmt.Println(server.proposer.V != otherServer.proposer.V)
	// fmt.Println(server.proposer.N_a_max != otherServer.proposer.N_a_max)
	// fmt.Println(server.proposer.Phase != otherServer.proposer.Phase)
	// fmt.Println(server.proposer.InitialValue != otherServer.proposer.InitialValue)

	if server.proposer.N != otherServer.proposer.N || server.proposer.V != otherServer.proposer.V ||
		server.proposer.N_a_max != otherServer.proposer.N_a_max || server.proposer.Phase != otherServer.proposer.Phase ||
		server.proposer.InitialValue != otherServer.proposer.InitialValue ||
		server.proposer.SuccessCount != otherServer.proposer.SuccessCount ||
		server.proposer.ResponseCount != otherServer.proposer.ResponseCount {
		return false
	}

	for i, response := range server.proposer.Responses {
		if response != otherServer.proposer.Responses[i] {
			return false
		}
	}

	return true
}

func (server *Server) Address() base.Address {
	return server.peers[server.me]
}
