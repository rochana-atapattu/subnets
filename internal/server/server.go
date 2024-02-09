package server

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rochana-atapattu/subnets/internal/subnet"
	"github.com/rochana-atapattu/subnets/internal/types"
	"github.com/rochana-atapattu/subnets/internal/view"
)

var root *subnet.Subnet

type Server struct {
	listenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}
func (s *Server) Run() {
	app := echo.New()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(c.Request().Context(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	app.Use(middleware.Recover())

	app.Static("/static", "static")
	app.GET("/", s.handleSubnetPage)
	app.GET("/reset", s.handleReset)
	app.POST("/calculate", s.handleSubnetCalculation)
	app.POST("/divide", s.handleSubnetDivision)
	app.POST("/join", s.handleSubnetJoining)

	app.Start(":8080")
}


func (s Server) handleSubnetPage(c echo.Context) error {
	si := view.SubnetIndex("subnet", view.Subnet())
	return render(c, si)

}
func (s Server) handleReset(c echo.Context) error {
	root = nil
	return c.Redirect(302, "/")
}

func (s Server) handleSubnetCalculation(c echo.Context) error {
	var req types.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	root = &subnet.Subnet{
		Address: subnet.InetAton(req.Network), 
		MaskLen: uint32(req.Netbits), 
		Parent: &subnet.Subnet{
			Address: subnet.InetAton(req.Network), 
			MaskLen: uint32(req.Netbits),
		},
	}
	tc := viewRows(root)
	return render(c, tc)
}

func (s Server) handleSubnetDivision(c echo.Context) error {
	var req types.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	slog.InfoContext(c.Request().Context(), "Dividing subnet: %s with mask length: %s", req.Network, fmt.Sprint(req.Netbits))
	// check if request is valid
	if req.Network == "" || req.Netbits == 0 {
		return c.String(400, "Invalid request")
	}
	address := subnet.InetAton(req.Network)
	maskLen := uint32(req.Netbits)
	root.Find(address, maskLen).Divide()
	tc := viewRows(root)
	return render(c, tc)
}

func (s Server) handleSubnetJoining(c echo.Context) error {
	var req types.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	slog.InfoContext(c.Request().Context(), "Joining subnet: %s with mask length: %s", req.Network, fmt.Sprint(req.Netbits))
	// check if request is valid
	if req.Network == "" || req.Netbits == 0 {
		return c.String(400, "Invalid request")
	}
	address := subnet.InetAton(req.Network)
	maskLen := uint32(req.Netbits)
	root.Find(address, maskLen).Join()
	tc := viewRows(root)
	return render(c, tc)
}
