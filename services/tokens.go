package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/dao/postgres"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

func (s *ServiceFacade) UpdateTokens() {
	tn := time.Now()
	tokenIdents, err := s.node.GetESDTs()
	if err != nil {
		log.Error("UpdateTokens: node.GetESDTs: %s", err.Error())
		return
	}
	for _, tokenIdent := range tokenIdents {
		props, err := s.node.GetESDTProperties(tokenIdent)
		if err != nil {
			log.Error("UpdateTokens: node.GetESDTProperties: %s", err.Error())
		}
		switch props.Type {
		case dmodels.FungibleESDT, dmodels.MetaESDT:
			err = s.updateFungibleToken(tokenIdent, props)
			if err != nil {
				log.Error("UpdateTokens: updateFungibleToken: %s", err.Error())
			}
		case dmodels.NonFungibleESDT, dmodels.SemiFungibleESDT:
			err = s.updateNonFungibleToken(tokenIdent, props)
			if err != nil {
				log.Error("UpdateTokens: updateNonFungibleToken: %s", err.Error())
			}
		default:
			log.Warn("UpdateTokens: unknown type: %s", props.Type)
		}
	}
	log.Debug("UpdateTokens: complete. duration: %s", time.Now().Sub(tn))
}

func (s *ServiceFacade) updateFungibleToken(tokenIdent string, props node.ESDTProperties) error {
	roles, err := s.node.GetESDTAllAddressesAndRoles(tokenIdent)
	if err != nil {
		return fmt.Errorf("node.GetESDTAllAddressesAndRoles: %s", err.Error())
	}
	supply, err := s.node.GetESDTSupply(tokenIdent)
	if err != nil {
		return fmt.Errorf("node.GetESDTAllAddressesAndRoles: %s", err.Error())
	}
	if props.Decimals > 0 {
		supply = supply.Div(decimal.New(1, int32(props.Decimals)))
	}
	rolesJSON, _ := json.Marshal(roles)
	propsJSON, _ := json.Marshal(props)
	splittedIdent := strings.Split(tokenIdent, "-")
	token := dmodels.Token{
		Identity:   tokenIdent,
		Name:       splittedIdent[0],
		Type:       props.Type,
		Owner:      props.Owner,
		Supply:     supply,
		Decimals:   uint64(props.Decimals),
		Properties: propsJSON,
		Roles:      rolesJSON,
	}
	_, err = s.dao.GetToken(tokenIdent)
	if err != nil {
		if err.Error() == postgres.NoRowsError {
			err = s.dao.CreateToken(token)
			if err != nil {
				return fmt.Errorf("dao.CreateToken: %s", err.Error())
			}
		}
	} else {
		err = s.dao.UpdateToken(token)
		if err != nil {
			return fmt.Errorf("dao.UpdateToken: %s", err.Error())
		}
	}
	return nil
}

func (s *ServiceFacade) updateNonFungibleToken(tokenIdent string, props node.ESDTProperties) error {
	propsJSON, _ := json.Marshal(props)
	tokenInfo, err := s.dao.GetTokenInfo(tokenIdent)
	if err != nil {
		return fmt.Errorf("dao.GetTokenInfo: %s", err.Error())
	}
	name := tokenInfo.Name
	if name == "" {
		s := strings.Split(tokenIdent, "-")
		name = s[0]
	}
	collection := dmodels.NFTCollection{
		Identity:   tokenIdent,
		Name:       name,
		Owner:      props.Owner,
		Type:       props.Type,
		Properties: propsJSON,
		CreatedAt:  time.Unix(int64(tokenInfo.Timestamp), 0),
	}
	_, err = s.dao.GetNFTCollection(tokenIdent)
	if err != nil {
		if err.Error() == postgres.NoRowsError {
			err = s.dao.CreateNFTCollection(collection)
			if err != nil {
				return fmt.Errorf("dao.CreateToken: %s", err.Error())
			}
		}
	} else {
		err = s.dao.UpdateNFTCollection(collection)
		if err != nil {
			return fmt.Errorf("dao.UpdateToken: %s", err.Error())
		}
	}
	return nil
}

func (s *ServiceFacade) GetNFT(id string) (sNFT smodels.NFT, err error) {
	nft, err := s.dao.GetTokenInfo(id)
	if err != nil {
		return sNFT, fmt.Errorf("dao.GetTokenInfo: %s", err.Error())
	}
	if nft.Data == nil {
		return sNFT, fmt.Errorf("empty data in nft token %s", nft.Identifier)
	}
	return toNFTSModel(nft), nil
}

func (s *ServiceFacade) GetNFTs(filter filters.NFTTokens) (pagination smodels.Pagination, err error) {
	tokens, err := s.dao.GetNFTTokens(filter)
	if err != nil {
		return pagination, fmt.Errorf("dao.GetCollectionNFTs: %s", err.Error())
	}
	nfts := make([]smodels.NFT, len(tokens))
	for i, t := range tokens {
		if t.Data == nil {
			return pagination, fmt.Errorf("empty data in nft token %s", t.Identifier)
		}
		nfts[i] = toNFTSModel(t)
	}
	total, err := s.dao.GetNFTTokensCount(filter)
	if err != nil {
		return pagination, fmt.Errorf("dao.GetNFTTokensCount: %s", err.Error())
	}
	return smodels.Pagination{
		Items: nfts,
		Count: total,
	}, nil
}

func toNFTSModel(nft data.TokenInfo) smodels.NFT {
	assets := make([]string, len(nft.Data.URIs))
	for i, u := range nft.Data.URIs {
		assets[i] = base64.StdEncoding.EncodeToString(u)
	}
	assetsJSON, _ := json.Marshal(assets)
	name := nft.Name
	if nft.Data != nil {
		name = nft.Data.Name
	}
	return smodels.NFT{
		Name:       name,
		Identity:   nft.Identifier,
		Owner:      nft.CurrentOwner,
		Creator:    nft.Data.Creator,
		Collection: nft.Token,
		Type:       nft.Type,
		Minted:     smodels.NewTime(time.Unix(int64(nft.Timestamp), 0)),
		Royalties:  decimal.New(int64(nft.Data.Royalties), -2),
		Assets:     assetsJSON,
	}
}

func (s *ServiceFacade) GetNFTCollection(id string) (collection smodels.NFTCollection, err error) {
	col, err := s.dao.GetNFTCollection(id)
	if err != nil {
		return collection, fmt.Errorf("dao.GetNFTCollection: %s", err.Error())
	}
	return toNFTCollectionSModel(col), nil
}

func (s *ServiceFacade) GetNFTCollections(filter filters.NFTCollections) (pagination smodels.Pagination, err error) {
	collections, err := s.dao.GetNFTCollections(filter)
	if err != nil {
		return pagination, fmt.Errorf("dao.GetNFTCollections: %s", err.Error())
	}
	total, err := s.dao.GetNFTCollectionsTotal(filter)
	if err != nil {
		return pagination, fmt.Errorf("dao.GetNFTCollectionsTotal: %s", err.Error())
	}
	items := make([]smodels.NFTCollection, len(collections))
	for i, collection := range collections {
		items[i] = toNFTCollectionSModel(collection)
	}
	return smodels.Pagination{
		Items: items,
		Count: total,
	}, nil
}

func toNFTCollectionSModel(collection dmodels.NFTCollection) smodels.NFTCollection {
	return smodels.NFTCollection{
		Name:       collection.Name,
		Identity:   collection.Identity,
		Owner:      collection.Owner,
		Type:       collection.Type,
		Properties: collection.Properties,
		CreatedAt:  smodels.NewTime(collection.CreatedAt),
	}
}

func (s *ServiceFacade) GetToken(id string) (token smodels.Token, err error) {
	t, err := s.dao.GetToken(id)
	if err != nil {
		return token, fmt.Errorf("dao.GetToken: %s", err.Error())
	}
	return toTokenSModel(t), nil
}

func (s *ServiceFacade) GetTokens(filter filters.Tokens) (pagination smodels.Pagination, err error) {
	tokens, err := s.dao.GetTokens(filter)
	if err != nil {
		return pagination, fmt.Errorf("dao.GetTokens: %s", err.Error())
	}
	total, err := s.dao.GetTokensCount(filter)
	if err != nil {
		return pagination, fmt.Errorf("dao.GetTokensCount: %s", err.Error())
	}
	items := make([]smodels.Token, len(tokens))
	for i, t := range tokens {
		items[i] = toTokenSModel(t)
	}
	return smodels.Pagination{
		Items: items,
		Count: total,
	}, nil
}

func toTokenSModel(token dmodels.Token) smodels.Token {
	return smodels.Token{
		Identity:   token.Identity,
		Name:       token.Name,
		Type:       token.Type,
		Owner:      token.Owner,
		Supply:     token.Supply,
		Decimals:   token.Decimals,
		Properties: token.Properties,
		Roles:      token.Roles,
	}
}
