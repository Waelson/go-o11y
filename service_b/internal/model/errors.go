package model

import "errors"

var InvalidCepError = errors.New("cep vazio")
var CepNotFoundError = errors.New("cep nao encontrado")
var InternalError = errors.New("ocorreu um erro interno")
