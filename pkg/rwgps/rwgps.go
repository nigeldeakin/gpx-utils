package rwgps

import (
	"fmt"
	"io"
	"net/http"
)

type ErrNotFound struct {
	RouteId int
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("RideWithGPS track %d not found", e.RouteId)
}

type ErrNotPublic struct {
	RouteId int
}

func (e *ErrNotPublic) Error() string {
	return fmt.Sprintf("RideWithGPS track %d is not public", e.RouteId)
}

func FetchTrack(routeId int) ([]byte, error) {

  url := fmt.Sprintf("https://ctccambridge.org.uk/getroutegpx?id=%d", routeId)
  
  // create custom HTTP client
  client := &http.Client{
    Transport: &http.Transport{},
  }

  // create HTTP request
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    return nil, fmt.Errorf("error creating HTTP request for %s: %v", url, err)
  }
 
  // set User-Agent header
  req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

  // make HTTP request
  response, err := client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("error getting %s: %v", url, err)
  }
  // close the response body
  defer response.Body.Close()

  if response.StatusCode != http.StatusOK {
    if response.StatusCode == http.StatusNotFound {
      return nil, &ErrNotFound{routeId}
    }
    if response.StatusCode == http.StatusForbidden {
      return nil, &ErrNotPublic{routeId}
    }
    return nil, fmt.Errorf("error retrieving route %d from %s, status is %s", routeId, url,response.Status)
  }

  // read the response body
  body, err := io.ReadAll(response.Body)
  if err != nil {
    return nil, fmt.Errorf("error reading response body for %s: %v", url, err)
  }

	return body, nil
}
