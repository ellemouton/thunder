package server

import (
        "context"
        blogs_db "github.com/ellemouton/thunder/elle/blogs/db"
        pb "github.com/ellemouton/thunder/elle/ellepb"
        "strconv"
        "strings"
)

var _ pb.ElleServer = (*Server)(nil)

type Server struct{
        b Backends
}

func New(b Backends) *Server {
        return &Server{b: b}
}

func (s Server) RequiresPayment(ctx context.Context, request *pb.RequiresPaymentRequest) (*pb.RequiresPaymentResponse, error) {
       if !strings.Contains(request.Path, "blog/view") {
               return &pb.RequiresPaymentResponse{
                       RequiresPayment: false,
               }, nil
       }

       id, err := strconv.Atoi(strings.TrimLeft(request.Path, "blog/view"))
       if err != nil {
               return &pb.RequiresPaymentResponse{
                       RequiresPayment: false,
               }, err
       }

       info, err := blogs_db.LookupInfo(ctx, s.b.GetDB(), int64(id))
       if err != nil {
               return &pb.RequiresPaymentResponse{
                       RequiresPayment: false,
               }, err
       }

       return &pb.RequiresPaymentResponse{
               RequiresPayment: info.Price != 0,
               Price:           info.Price,
               Id:              info.ID,
               Memo:            info.Description,
       }, nil
}

