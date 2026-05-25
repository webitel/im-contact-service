package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/webitel/webitel-go-kit/pkg/errors"

	impb "github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/service"
)

type ViaServer struct {
	impb.UnimplementedViasServer

	via service.ViaService
}

func newViaServer(via service.ViaService) *ViaServer {
	return &ViaServer{via: via}
}

func (viaServer *ViaServer) Create(ctx context.Context, req *impb.CreateViaRequest) (*impb.Via, error) {
	contactID := uuid.Nil

	if req.GetContactId() != "" {
		parsedContactID, err := uuid.Parse(req.GetContactId())
		if err != nil {
			return nil, errors.InvalidArgument("contact id has invalid format", errors.WithCause(err), errors.WithID("grpc.via.create"))
		}

		contactID = parsedContactID
	}

	via := &model.CreateViaCommunicationCommand{
		ContactID:     contactID,
		Via:           req.GetVia(),
		Disable:       req.GetDisable(),
		DisableReason: req.DisableReason,
		Metadata:      req.GetMetadata().AsMap(),
		Iss:           req.GetIss(),
		Sub:           req.GetSub(),
	}

	created, err := viaServer.via.Create(ctx, via)
	if err != nil {
		return nil, err
	}

	result, err := convertDomainToProto(created)
	if err != nil {
		return nil, errors.Wrap(err, errors.WithID("grpc.via.create"))
	}

	return result, nil
}

func (viaServer *ViaServer) Update(ctx context.Context, req *impb.UpdateViaRequest) (*impb.Via, error) {
	contactID, err := uuid.Parse(req.GetContactId())
	if err != nil {
		return nil, errors.InvalidArgument("contact id has invalid uuid format", errors.WithCause(err), errors.WithID("grpc.via.update"))
	}

	update := &model.ViaCommunication{
		ContactID:     contactID,
		Via:           req.GetVia(),
		Disable:       req.GetDisable(),
		DisableReason: req.DisableReason,
		Metadata:      req.GetMetadata().AsMap(),
	}

	updated, err := viaServer.via.Update(ctx, update)
	if err != nil {
		return nil, err
	}

	response, err := convertDomainToProto(updated)
	if err != nil {
		return nil, errors.Wrap(err, errors.WithID("grpc.via.update"))
	}

	return response, nil
}

func (viaServer *ViaServer) PartialUpdate(ctx context.Context, req *impb.PartialUpdateViaRequest) (*impb.Via, error) {
	var contactID uuid.UUID

	if req.GetUpdate().GetContactId() != "" {
		reid, err := uuid.Parse(req.GetUpdate().GetContactId())
		if err != nil {
			return nil, errors.InvalidArgument("contact id has invalid uuid format", errors.WithCause(err), errors.WithID("grpc.via.partial_update"))
		}

		contactID = reid
	}

	partialUpdateCmd := &model.CommunicationViaPartialUpdateCmd{
		ViaCommunication: model.ViaCommunication{
			ContactID:     contactID,
			Via:           req.GetUpdate().GetVia(),
			Disable:       req.GetUpdate().GetDisable(),
			DisableReason: req.GetUpdate().DisableReason,
			Metadata:      req.GetUpdate().GetMetadata().AsMap(),
		},
		Fields: req.GetFieldMask().GetPaths(),
	}

	updated, err := viaServer.via.PartialUpdate(ctx, partialUpdateCmd)
	if err != nil {
		return nil, err
	}

	response, err := convertDomainToProto(updated)
	if err != nil {
		return nil, errors.Wrap(err, errors.WithID("grpc.via.partial_update"))
	}

	return response, nil
}

func (viaServer *ViaServer) Search(ctx context.Context, req *impb.SearchViaRequest) (*impb.SearchViaResponse, error) {
	contactIDs, err := convertStringsToUUID(req.GetContactIds())
	if err != nil {
		return nil, errors.Wrap(err, errors.WithID("grpc.via.search"))
	}

	page, size := ParsePagination(int32(req.GetPage()), int32(req.GetSize()))

	filter := &model.SearchViaCommunicationsFilter{
		Sort:       req.GetSort(),
		Limit:      int(size),
		Page:       int(page),
		Fields:     req.GetFields(),
		ContactIDs: contactIDs,
		Disabled:   req.Disabled,
		Vias:       req.GetVias(),
	}

	vias, err := viaServer.via.Search(ctx, filter)
	if err != nil {
		return nil, err
	}

	protoVias, err := convertDomainViaListToProto(vias)
	if err != nil {
		return nil, err
	}

	parsed, next := ResolvePaging(int(size), protoVias)
	response := &impb.SearchViaResponse{
		Items: parsed,
		Next:  next,
		Page:  page,
	}

	return response, nil
}

func convertStringsToUUID(idStrs []string) (uuid.UUIDs, error) {
	idStrsLen := len(idStrs)
	if idStrsLen == 0 {
		return nil, nil
	}

	err := (error)(nil)

	ids := make([]uuid.UUID, idStrsLen)
	for i, s := range idStrs {
		if ids[i], err = uuid.Parse(s); err != nil {
			return nil, errors.InvalidArgument("converting contact string id to uuid", errors.WithCause(err), errors.WithID("grpc.via.search"), errors.WithValue("index", i))
		}
	}

	return ids, nil
}

func convertDomainViaListToProto(vias []*model.ViaCommunication) ([]*impb.Via, error) {
	inputLen := len(vias)
	err := (error)(nil)
	proto := make([]*impb.Via, inputLen)

	for i, via := range vias {
		if proto[i], err = convertDomainToProto(via); err != nil {
			return nil, err
		}
	}

	return proto, nil
}

func convertDomainToProto(via *model.ViaCommunication) (*impb.Via, error) {
	if via == nil {
		return new(impb.Via), nil
	}

	metadata, err := structpb.NewStruct(via.Metadata)
	if err != nil {
		return nil, errors.Internal("converting via metadata to proto struct", errors.WithCause(err), errors.WithID("grpc.via.convert_domain_to_proto"))
	}

	return &impb.Via{
		ContactId:     via.ContactID.String(),
		Via:           via.Via,
		Disable:       via.Disable,
		DisableReason: via.DisableReason,
		CreatedAt:     via.CreatedAtUTCUnix(),
		UpdatedAt:     via.UpdatedAtUTCUnix(),
		Metadata:      metadata,
	}, nil
}
