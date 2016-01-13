#!/bin/bash
#
# this script processes all .html files in the directory, inlining css and
# striping script/noscript tags.
#
# make sure 'critical' is installed:
# sudo npm install -g critical
#
# make sure goquery is installed:
# go get -u github.com/PuerkitoBio/goquery

SRC=$(realpath $(cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd ))
PKGNAME=$(basename "$SRC")

GOSTRIPDIR=$(mktemp -d)
GOSTRIPPATH=$GOSTRIPDIR/gostrip.go

# gostrip src
echo 'package main

import (
  "fmt"
  "log"
  "os"

  "github.com/PuerkitoBio/goquery"
)

func main() {
  doc, err := goquery.NewDocumentFromReader(os.Stdin)
  if err != nil {
    log.Fatalln(err)
  }

  doc.Find("script").Each(func(i int, s *goquery.Selection) {
    s.Remove()
  })

  doc.Find("noscript").Each(func(i int, s *goquery.Selection) {
    s.Remove()
  })

  html, err := doc.Html()
  if err != nil {
    log.Fatalln(err)
  }

  fmt.Println(html)
}' > $GOSTRIPPATH

# go header
GODATA=$(cat <<ENDHDR
// Package $PKGNAME is a collection of static template data.
package $PKGNAME

//go:generate ./gen.sh
ENDHDR
)

NL=$'\n'

# loop over files, process, and add to GODATA as a go const
for i in $SRC/*.html; do
  echo -n "processing $i ..."
  FN=$(basename $i .html)
  GODATA+="${NL}${NL}// ${FN} is the contents of ${FN}.html${NL}const ${FN} = \`$(critical $i --inline --extract --minify |perl -pe 's/<script[^>]+>.*<\/script>//igs' |go run $GOSTRIPPATH)\`"
  echo "done."
done

# remove temporary go dir
rm -rf $GOSTRIPDIR

# write go file
cat <<<"$GODATA" |goimports|gofmt > "${SRC}/${PKGNAME}.go"
