package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/morningconsult/go-elasticsearch-alerts/command/alert"
	"golang.org/x/xerrors"
)

var _ alert.Method = (*AlertMethod)(nil)

type AlertMethodConfig struct {
	Corpid  string            `mapstructure:"corpid" json:"corpid"`
	Secret  string            `mapstructure:"secret" json:"secret"`
	ToUser  string            `mapstructure:"touser" json:"touser"`
	ToParty string            `mapstructure:"toparty" json:"toparty"`
	ToTag   string            `mapstructure:"totag" json:"totag"`
	Agentid int               `mapstructure:"agentid" json:"agentid"`
	Text    map[string]string `mapstructure:"text" json:"text"`
	MsgType string            `mapstructure:"msgtype" json:"msgtype"`
	Safe    string            `mapstructure:"safe" json:"safe"`
}

// AlertMethod implements the alert.Method interface
// for writing new alerts to email.
type AlertMethod struct {
	Corpid  string            `json:"corpid"`
	Secret  string            `json:"secret"`
	Touser  string            `json:"touser"`
	Toparty string            `json:"toparty"`
	Totag   string            `json:"totag"`
	Agentid int               `json:"agentid"`
	Msgtype string            `json:"msgtype"`
	Text    map[string]string `json:"text"`
}

// NewAlertMethod creates a new *AlertMethod or a
// non-nil error if there was an error.
func NewAlertMethod(config *AlertMethodConfig) (alert.Method, error) {
	if config == nil {
		return nil, xerrors.New("no config provided")
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return &AlertMethod{
		Corpid:  config.Corpid,
		Secret:  config.Secret,
		Touser:  config.ToUser,
		Toparty: config.ToParty,
		Totag:   config.ToTag,
		Agentid: config.Agentid,
		Msgtype: config.MsgType,
	}, nil
}

func validateConfig(config *AlertMethodConfig) error {
	var allErrors *multierror.Error
	if config == nil {
		allErrors = multierror.Append(xerrors.New("no config provided"))
	}

	if config.Corpid == "" {
		allErrors = multierror.Append(xerrors.New("no corpid provided"))
	}
	if config.Secret == "" {
		allErrors = multierror.Append(xerrors.New("no secret provided"))
	}
	if config.ToUser == "" {
		allErrors = multierror.Append(xerrors.New("no touser provided"))
	}
	if config.Agentid == 0 {
		allErrors = multierror.Append(xerrors.New("no agentid provided"))
	}
	return allErrors.ErrorOrNil()
}

func (w *AlertMethod) Write(ctx context.Context, rule string, records []*alert.Record) error {
	body, err := w.BuildMessage(rule, records)
	if err != nil {
		return xerrors.Errorf("error creating WeChat message: %v", err)
	}
	w.Text = map[string]string{"content": body}
	return w.SendWeChat()
}

// BuildMessage creates an email message from the provided
// records. It will return a non-nil error if an error occurs.
func (w *AlertMethod) BuildMessage(rule string, records []*alert.Record) (string, error) {
	msg := fmt.Sprintf("%s\n", rule)
	for _, record := range records {
		text := record.Text
		if record.Fields != nil {
			var fields = ""
			for _, field := range record.Fields {
				fields += fmt.Sprintf("\n%s: %d", field.Key, field.Count)
			}
			text += fields
		}
		msg += text + "\n"
	}
	return msg, nil
}

func (w *AlertMethod) SendWeChat() error {
	token, err := Get_token(w.Corpid, w.Secret)
	if err != nil {
		return err
	}
	buf, err := json.Marshal(w)
	if err != nil {
		return err
	}
	return Send_msg(token.Access_token, buf)
}
