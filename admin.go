package broker

import (
	"errors"
	rh "github.com/michaelklishin/rabbit-hole"
	"net/http"
)

const (
	defaultUri      = "http://localhost:15672"
	defaultUser     = "guest"
	defaultPassword = "guest"
)

var (
	errorAlreadyExists      = errors.New("Admin: Entity already exists.")
	errorUnexpectedResponse = errors.New("Admin: Unexpected response received.")
)

var defaultOptions = options{
	uri:      defaultUri,
	user:     defaultUser,
	password: defaultPassword,
}

type options struct {
	uri, user, password string
}

type rabbitAdmin struct {
	client *rh.Client
}

func newRabbitAdmin(opt options) (*rabbitAdmin, error) {
	client, err := rh.NewClient(opt.uri, opt.username, opt.password)
	if err != nil {
		return nil, err
	}
	return &rabbitAdmin{client}, nil
}

func (a *rabbitAdmin) isVhost(username string) (bool, error) {
	info, err := a.client.GetVHost(username)
	if info != nil {
		return true, nil
	} else if err.Error() == "not found" {
		return false, nil
	}
	return false, err
}

func (a *rabbitAdmin) createVhost(vhostname string, tracing bool) error {
	if found, err := a.isVhost(vhostname); err != nil {
		return err
	} else if found {
		return errorAlreadyExists
	}

	settings := rh.VhostSettings{tracing}
	resp, err := a.client.PutVhost(vhostname, settings)
	if err != nil {
		return err
	}
	return checkResponseAndClose(resp)
}

func (a *rabbitAdmin) deleteVhost(vhostname string) error {
	resp, err := a.client.DeleteVhost(vhostname)
	if err != nil {
		return err
	}
	return checkResponseAndClose(resp)
}

func (a *rabbitAdmin) isUser(username string) (bool, error) {
	info, err := a.client.GetUser(username)
	if info != nil {
		return true, nil
	} else if err.Error() == "not found" {
		return false, nil
	}
	return false, err
}

func (a *rabbitAdmin) createUser(username, password string) error {
	if found, err := a.isUser(username); err != nil {
		return err
	} else if found {
		return errorAlreadyExists
	}

	settings := rh.UserSettings{
		Name:     username,
		Password: password,
		Tags:     "management",
	}
	resp, err := a.client.PutUser(username, settings)
	if err != nil {
		return err
	}
	return checkResponseAndClose(resp)
}

func (a *rabbitAdmin) deleteUser(username string) error {
	resp, err := a.client.DeleteUser(username)
	if err != nil {
		return err
	}
	return checkResponseAndClose(resp)
}

func (a *rabbitAdmin) grantAllPermissionsIn(username, vhostname string) {
	unlimited := rh.Permissions{".*", ".*", ".*"}
	resp, err := a.client.UpdatePermissionsIn(vhostname, username, unlimited)
	if err != nil {
		return err
	}
	return checkResponseAndClose(resp)
}

func checkResponseAndClose(resp *http.Response) error {
	defer resp.Close()
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return nil
	}
	// TODO: Log the actual response status
	return errorUnexpectedResponse
}
