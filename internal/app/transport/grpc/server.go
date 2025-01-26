package grpc

import (
	"context"
	"fmt"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	_ "google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"songs/internal/app/domain"
	pb "songs/internal/app/proto"
	"songs/internal/app/service"
	"strconv"
	"time"
)

type Server struct {
	pb.UnimplementedSongServiceServer
	songService *service.SongService
	addr        string
}

func NewServer(addr string, songService *service.SongService) *Server {
	return &Server{
		songService: songService,
		addr:        addr,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.addr, err)
	}

	grpcServer := googlegrpc.NewServer()
	pb.RegisterSongServiceServer(grpcServer, s)

	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on %s", s.addr)
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (s *Server) Stop() {
	//pass TODO: add method
}

func (s *Server) GetSong(ctx context.Context, req *pb.GetSongRequest) (*pb.GetSongResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "song ID is required")
	}

	songID, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid song ID format")
	}

	song, err := s.songService.GetSong(ctx, songID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get song")
	}

	return &pb.GetSongResponse{
		Song: &pb.Song{
			Id:          strconv.Itoa(song.ID),
			Group:       strconv.Itoa(song.GroupID),
			Name:        song.Title,
			ReleaseDate: song.ReleaseDate.Format("2006-01-02"),
			Text:        song.Text,
			Link:        song.Link,
		},
	}, nil
}

func (s *Server) ListSongs(ctx context.Context, req *pb.ListSongsRequest) (*pb.ListSongsResponse, error) {
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	filters := make(map[string]string)
	if req.Group != "" {
		filters["group"] = req.Group
	}
	if req.Song != "" {
		filters["song"] = req.Song
	}
	if req.ReleaseDate != "" {
		filters["release_date"] = req.ReleaseDate
	}
	if req.Text != "" {
		filters["text"] = req.Text
	}
	if req.Link != "" {
		filters["link"] = req.Link
	}

	songs, total, err := s.songService.GetSongs(ctx, filters, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list songs")
	}

	var pbSongs []*pb.Song
	for _, song := range songs {
		pbSongs = append(pbSongs, &pb.Song{
			Id:          strconv.Itoa(song.ID),
			Group:       strconv.Itoa(song.GroupID),
			Name:        song.Title,
			ReleaseDate: song.ReleaseDate.Format("2006-01-02"),
			Text:        song.Text,
			Link:        song.Link,
		})
	}

	totalInt := int(total)
	pages := (totalInt + int(req.PageSize) - 1) / int(req.PageSize)

	return &pb.ListSongsResponse{
		Songs: pbSongs,
		Total: total,
		Page:  req.Page,
		Pages: int32(pages),
	}, nil
}

func (s *Server) CreateSong(ctx context.Context, req *pb.CreateSongRequest) (*pb.CreateSongResponse, error) {
	groupID, err := strconv.Atoi(req.Group)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid group ID format")
	}

	releaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid release date format")
	}

	song := &domain.Song{
		GroupID:     groupID,
		Title:       req.Name,
		ReleaseDate: releaseDate,
		Text:        req.Text,
		Link:        req.Link,
	}

	createdSong, err := s.songService.CreateSong(ctx, song)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create song")
	}

	return &pb.CreateSongResponse{
		Song: &pb.Song{
			Id:          strconv.Itoa(createdSong.ID),
			Group:       strconv.Itoa(createdSong.GroupID),
			Name:        createdSong.Title,
			ReleaseDate: createdSong.ReleaseDate.Format("2006-01-02"),
			Text:        createdSong.Text,
			Link:        createdSong.Link,
		},
	}, nil
}

func (s *Server) UpdateSong(ctx context.Context, req *pb.UpdateSongRequest) (*pb.UpdateSongResponse, error) {
	songID, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid song ID format")
	}

	groupID, err := strconv.Atoi(req.Group)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid group ID format")
	}

	releaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid release date format")
	}

	song := &domain.Song{
		ID:          songID,
		GroupID:     groupID,
		Title:       req.Name,
		ReleaseDate: releaseDate,
		Text:        req.Text,
		Link:        req.Link,
	}

	updatedSong, err := s.songService.UpdateSong(ctx, songID, song)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update song")
	}

	return &pb.UpdateSongResponse{
		Song: &pb.Song{
			Id:          strconv.Itoa(updatedSong.ID),
			Group:       strconv.Itoa(updatedSong.GroupID),
			Name:        updatedSong.Title,
			ReleaseDate: updatedSong.ReleaseDate.Format("2006-01-02"),
			Text:        updatedSong.Text,
			Link:        updatedSong.Link,
		},
	}, nil
}

func (s *Server) DeleteSong(ctx context.Context, req *pb.DeleteSongRequest) (*pb.DeleteSongResponse, error) {
	songID, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid song ID format")
	}

	err = s.songService.DeleteSong(ctx, songID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete song")
	}

	return &pb.DeleteSongResponse{
		Success: true,
	}, nil
}
