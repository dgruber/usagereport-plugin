#!/bin/sh

go test -v 
go test github.com/dgruber/usagereport-plugin/apihelper -v
go test github.com/dgruber/usagereport-plugin/models -v
