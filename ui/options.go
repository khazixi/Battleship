package ui

import (
	"encoding/gob"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khazixi/Battelship/util"
)

type Options struct {
	Room   int
	Active bool
	msgch  chan util.Message
	conn   net.Conn
}

func MakeOption(conn net.Conn) Options {
	return Options{
		conn:  conn,
		msgch: make(chan util.Message),
	}
}

func (o Options) Listen(msg chan util.Message) tea.Cmd {
	return func() tea.Msg {
		decoder := gob.NewDecoder(o.conn)
		for {
			message, err := util.MessageDecoder(decoder)
			if err != nil {
				msg <- util.ErrorMessage{Err: err}
				continue
			}
			msg <- message
		}
	}
}

func (o Options) Process(msg chan util.Message) tea.Cmd {
	return func() tea.Msg {
		return <-msg
	}
}

func (o *Options) Send(action util.Action) tea.Cmd {
	return func() tea.Msg {
		encoder := gob.NewEncoder(o.conn)
		util.ActionEncoder(encoder, action)
		return 0
	}
}

func (o Options) Connect(ip string) tea.Cmd {
	return func() tea.Msg {
		conn, err := net.Dial("tcp", ip)
		if err != nil {
			return err
		}
		return conn
	}
}
