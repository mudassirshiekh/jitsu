package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jitsucom/jitsu/server/logging"
	"github.com/jitsucom/jitsu/server/middleware"
	"github.com/jitsucom/jitsu/server/schema"
	"github.com/jitsucom/jitsu/server/templates"
	"net/http"
)

//EvaluateTemplateRequest is a request dto for testing text/template expressions
type EvaluateTemplateRequest struct {
	Object     map[string]interface{} `json:"object,omitempty"`
	Expression string                 `json:"expression,omitempty"`
	Reformat   bool                   `json:"reformat,omitempty"`
}

//EvaluateTemplateResponse is a response dto for testing text/template expressions
type EvaluateTemplateResponse struct {
	Result string `json:"result"`
	Error string `json:"message"`
	Format string `json:"format"`
}

//Validate returns err if invalid
func (etr *EvaluateTemplateRequest) Validate() error {
	if etr.Object == nil {
		return errors.New("'object' is required field")
	}

	if etr.Expression == "" {
		return errors.New("'expression' is required field")
	}

	return nil
}

//EventTemplateHandler is a handler for testing text/template expression with income object
func EventTemplateHandler(c *gin.Context) {
	req := &EvaluateTemplateRequest{}
	if err := c.BindJSON(req); err != nil {
		logging.Errorf("Error parsing evaluate template body: %v", err)
		c.JSON(http.StatusBadRequest, middleware.ErrResponse("Failed to parse body", err))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, middleware.ErrResponse(err.Error(), nil))
		return
	}

	var result string
	var format string
	var err error

	if req.Reformat {
		result, format, err = evaluateReformatted(req)
	} else {
		result, format, err = evaluate(req)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, EvaluateTemplateResponse{Result: result, Format: format, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, EvaluateTemplateResponse{Result: result, Format: format})
}

func evaluate(req *EvaluateTemplateRequest) (result string, format string, err error) {
	//panic handler
	defer func() {
		if r := recover(); r != nil {
			result = ""
			err = fmt.Errorf("Error: %v", r)
		}
	}()

	tmpl, err := templates.SmartParse("template evaluating", req.Expression, templates.JSONSerializeFuncs)
	if err != nil {
		return "", "", fmt.Errorf("error parsing template: %v", err)
	}
	resultObject, err:= tmpl.ProcessEvent(req.Object)
	if err != nil {
		return "", tmpl.Format(), fmt.Errorf("error executing template: %v", err)
	}
	jsonBytes, err := templates.ToJSON(resultObject)
	if err != nil {
		return "", tmpl.Format(), err
	}
	result = string(jsonBytes)
	format = tmpl.Format()
	return
}

func evaluateReformatted(req *EvaluateTemplateRequest) (string, string, error) {
	tableNameExtractor, err := schema.NewTableNameExtractor(req.Expression)
	if err != nil {
		return "", "", err
	}
	res, err := tableNameExtractor.Extract(req.Object)
	return res, tableNameExtractor.Format(), err
}
