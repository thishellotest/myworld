package config_vbc

import "vbc/configs"

const User_Yannan_gid = "6159272000000453001"
const User_Edward_gid = "6159272000000453640"
const User_Dev_gid = "6159272000000453669" // glliao
const User_Victoria_gid = "6159272000001027129"
const User_Lili_gid = "6159272000001027094"
const User_Sharika_gid = "ff28b522130446158751d5f6c273163b"
const User_Elissa_gid = "b796f7957fbe4dcba6e0f11d3545c2ef"

func GetUserVictoriaGid() string {
	if configs.IsProd() {
		return User_Victoria_gid
	}
	return User_Dev_gid
}

func GetUserEdwardGid() string {
	if configs.IsProd() {
		return User_Edward_gid
	}
	return User_Dev_gid
}

func GetUserElissaGid() string {
	if configs.IsProd() {
		return User_Elissa_gid
	}
	return User_Dev_gid
}

func GetUserSharikaGid() string {
	if configs.IsProd() {
		return User_Sharika_gid
	}
	return User_Dev_gid
}

const Role_VSTeam_gid = "ba8b82363dd646e0856c62b00402f978"

type State struct {
	FullName  string
	ShortName string
}

type ListState []State

func (c ListState) FullNameByShort(shortName string) string {
	for k, v := range c {
		if v.ShortName == shortName {
			return c[k].FullName
		}
	}
	return ""
}

// website没有的：
// {FullName: "American Samoa", ShortName: ""} {FullName: "Guam", ShortName: ""}
// {FullName: "Northern Mariana Islands", ShortName: ""} {FullName: "Puerto Rico", ShortName: ""}
// {FullName: "Trust Territories", ShortName: ""} {FullName: "Wyoming", ShortName: ""}
var StateConfigs = ListState{
	{FullName: "Alabama", ShortName: "AL"}, {FullName: "Alaska", ShortName: "AK"}, {FullName: "Arizona", ShortName: "AZ"},
	{FullName: "Arkansas", ShortName: "AR"}, {FullName: "American Samoa", ShortName: ""}, {FullName: "California", ShortName: "CA"},
	{FullName: "Colorado", ShortName: "CO"}, {FullName: "Connecticut", ShortName: "CT"}, {FullName: "Delaware", ShortName: "DE"},
	{FullName: "District of Columbia", ShortName: "DC"}, {FullName: "Florida", ShortName: "FL"}, {FullName: "Georgia", ShortName: "GA"},
	{FullName: "Guam", ShortName: ""}, {FullName: "Hawaii", ShortName: "HI"}, {FullName: "Idaho", ShortName: "ID"},
	{FullName: "Illinois", ShortName: "IL"}, {FullName: "Indiana", ShortName: "IN"}, {FullName: "Iowa", ShortName: "IA"},
	{FullName: "Kansas", ShortName: "KS"}, {FullName: "Kentucky", ShortName: "KY"}, {FullName: "Louisiana", ShortName: "LA"},
	{FullName: "Maine", ShortName: "ME"}, {FullName: "Maryland", ShortName: "MD"}, {FullName: "Massachusetts", ShortName: "MA"},
	{FullName: "Michigan", ShortName: "MI"}, {FullName: "Minnesota", ShortName: "MN"}, {FullName: "Mississippi", ShortName: "MS"},
	{FullName: "Missouri", ShortName: "MO"}, {FullName: "Montana", ShortName: "MT"}, {FullName: "Nebraska", ShortName: "NE"},
	{FullName: "Nevada", ShortName: "NV"}, {FullName: "New Hampshire", ShortName: "NH"}, {FullName: "New Jersey", ShortName: "NJ"},
	{FullName: "New Mexico", ShortName: "NM"}, {FullName: "New York", ShortName: "NY"}, {FullName: "North Carolina", ShortName: "NC"},
	{FullName: "North Dakota", ShortName: "ND"}, {FullName: "Northern Mariana Islands", ShortName: ""}, {FullName: "Ohio", ShortName: "OH"},
	{FullName: "Oklahoma", ShortName: "OK"}, {FullName: "Oregon", ShortName: "OR"}, {FullName: "Pennsylvania", ShortName: "PA"},
	{FullName: "Puerto Rico", ShortName: ""}, {FullName: "Rhode Island", ShortName: "RI"}, {FullName: "South Carolina", ShortName: "SC"},
	{FullName: "South Dakota", ShortName: "SD"}, {FullName: "Tennessee", ShortName: "TN"}, {FullName: "Texas", ShortName: "TX"},
	{FullName: "Trust Territories", ShortName: ""}, {FullName: "Utah", ShortName: "UT"}, {FullName: "Vermont", ShortName: "VT"},
	{FullName: "Virginia", ShortName: "VA"}, {FullName: "Virgin Islands", ShortName: ""}, {FullName: "Washington", ShortName: "WA"},
	{FullName: "West Virginia", ShortName: "WV"}, {FullName: "Wisconsin", ShortName: "WI"}, {FullName: "Wyoming", ShortName: ""},
	{FullName: "Outside United States", ShortName: "Outside United States"},
}
