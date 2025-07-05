package dto

import "mira/anima/datetime"

// Login Log List
type LogininforListResponse struct {
	InfoId        int               `json:"infoId"`
	UserName      string            `json:"userName"`
	Ipaddr        string            `json:"ipaddr"`
	LoginLocation string            `json:"loginLocation"`
	Browser       string            `json:"browser"`
	Os            string            `json:"os"`
	Status        string            `json:"status"`
	Msg           string            `json:"msg"`
	LoginTime     datetime.Datetime `json:"loginTime"`
}

// Login Log Export
type LogininforExportResponse struct {
	InfoId        int    `excel:"name:ID;"`
	UserName      string `excel:"name:User Account;"`
	Status        string `excel:"name:Login Status;replace:0_Success,1_Failure;"`
	Ipaddr        string `excel:"name:Login Address;"`
	LoginLocation string `excel:"name:Login Location;"`
	Browser       string `excel:"name:Browser;"`
	Os            string `excel:"name:Operating System;"`
	Msg           string `excel:"name:Message;"`
	LoginTime     string `excel:"name:Access Time;"`
}
