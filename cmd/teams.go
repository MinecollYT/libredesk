package main

import (
	"fmt"
	"strconv"

	"github.com/abhinavxd/artemis/internal/envelope"
	"github.com/abhinavxd/artemis/internal/team/models"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

func handleGetTeams(r *fastglue.Request) error {
	var (
		app        = r.Context.(*App)
	)
	teams, err := app.teamManager.GetAll()
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(teams)
}

func handleGetTeam(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
	)
	id, err := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	if err != nil || id == 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest,
			"Invalid team `id`.", nil, envelope.InputError)
	}
	team, err := app.teamManager.GetTeam(id)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(team)
}

func handleCreateTeam(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.Team{}
	)

	if _, err := fastglue.ScanArgs(r.RequestCtx.PostArgs(), &req, `json`); err != nil {
		app.lo.Error("error scanning args", "error", err)
		return envelope.NewError(envelope.InputError,
			fmt.Sprintf("Invalid request (%s)", err.Error()), nil)
	}
	err := app.teamManager.CreateTeam(req)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}

func handleUpdateTeam(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
		req = models.Team{}
	)
	id, err := strconv.Atoi(r.RequestCtx.UserValue("id").(string))
	if err != nil || id == 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest,
			"Invalid team `id`.", nil, envelope.InputError)
	}

	if _, err := fastglue.ScanArgs(r.RequestCtx.PostArgs(), &req, `json`); err != nil {
		return envelope.NewError(envelope.InputError,
			fmt.Sprintf("Invalid request (%s)", err.Error()), nil)
	}

	err = app.teamManager.UpdateTeam(id, req)
	if err != nil {
		return sendErrorEnvelope(r, err)
	}
	return r.SendEnvelope(true)
}
