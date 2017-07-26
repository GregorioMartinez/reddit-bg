package main

import (
	"errors"
	"math/rand"
	"net/url"
	"time"
)

type response struct {
	Data struct {
		Children []struct {
			Data postData `json:"data"`
			Kind string   `json:"kind"`
		} `json:"children"`
	} `json:"data"`
	Kind string `json:"kind"`
}

type postData struct {
	Title   string `json:"title"`
	Preview struct {
		Images []struct {
			ID          string `json:"id"`
			Resolutions []struct {
				Height float64 `json:"height"`
				URL    string  `json:"url"`
				Width  float64 `json:"width"`
			} `json:"resolutions"`
			Source struct {
				Height float64 `json:"height"`
				URL    string  `json:"url"`
				Width  float64 `json:"width"`
			} `json:"source"`
		} `json:"images"`
	} `json:"preview"`
}

func (r *response) selectPost(w float64, h float64, isRandom bool) (*postData, error) {
	if isRandom {
		random := rand.New(rand.NewSource(time.Now().Unix()))
		shuffled := random.Perm(len(r.Data.Children))
		for _, n := range shuffled {
			for _, v := range r.Data.Children[n].Data.Preview.Images {
				if v.Source.Height > h && v.Source.Width > w {
					return &r.Data.Children[n].Data, nil
				}
			}
		}
	} else {
		for _, n := range r.Data.Children {
			for _, v := range n.Data.Preview.Images {
				if v.Source.Height > h && v.Source.Width > w {
					return &n.Data, nil
				}
			}
		}
	}
	return nil, errors.New("Unable to find image that matches resolution")
}

func (p *postData) getImageURL() (*url.URL, error) {
	for _, v := range p.Preview.Images {
		u, err := url.Parse(v.Source.URL)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	return nil, errors.New("Unable to find source image")
}
