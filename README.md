# exp

![GitHub](https://img.shields.io/github/license/smallnest/exp) ![GitHub Action](https://github.com/smallnest/exp/actions/workflows/action.yaml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/exp)](https://goreportcard.com/report/github.com/smallnest/exp)  [![GoDoc](https://godoc.org/github.com/smallnest/exp?status.png)](http://godoc.org/github.com/smallnest/exp)  



Experimental packages not in std and golang.org/exp


## packages

- **sync**
  - generic sync.Map: modify sync.Map to support generic

- **container**
  - **heap**: generic heap
  - **list**: generic list
  - **ring**: generic ring
  - **skiplist**: generic skiplist based on [mauricegit/skiplist](https://github.com/mauricegit/skiplist)
  - **set**: discussion at https://github.com/golang/go/discussions/47331 
  - **maps**:
    - **OrderedMap**: an insert-order map. The main code is forked [wk8/go-ordered-map](https://github.com/wk8/go-ordered-map)
    - **BidiMap**: a bidirectional map. 

- **stat**
  - **win**: sliding window ![](stat/win/sliding.jpg)