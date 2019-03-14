#!/usr/bin/env bash

set -e

echo "==> Packaging binaries..."

source $(dirname $0)/version.sh

DIST=$(pwd)/dist/artifacts

mkdir -p $DIST/${VERSION} $DIST/latest

for i in build/bin/*; do
    if [ ! -e $i ]; then
        continue
    fi

    BASE=build/archive
    DIR=${BASE}/terraform-provider-rancher2-${VERSION}

    rm -rf $BASE
    mkdir -p $BASE $DIR

    EXT=
    if [[ $i =~ .*windows.* ]]; then
        EXT=.exe
    fi

    cp $i ${DIR}/terraform-provider-rancher2${EXT}

    arch=$(echo $i | cut -f2 -d_)
    mkdir -p $DIST/${VERSION}/binaries/$arch
    mkdir -p $DIST/latest/binaries/$arch
    cp $i $DIST/${VERSION}/binaries/$arch/terraform-provider-rancher2${EXT}
    if [ -z "${EXT}" ]; then
        gzip -c $i > $DIST/${VERSION}/binaries/$arch/terraform-provider-rancher2.gz
        xz -c $i > $DIST/${VERSION}/binaries/$arch/terraform-provider-rancher2.xz
    fi

    rm -rf $DIST/latest/binaries/$arch
    mkdir -p $DIST/latest/binaries
    cp -rf $DIST/${VERSION}/binaries/$arch $DIST/latest/binaries

    (
        cd $BASE
        NAME=$(basename $i | sed 's/_/-/g')
        if [ -z "$EXT" ]; then
            tar cvzf $DIST/${VERSION}/${NAME}-${VERSION}.tar.gz .
            cp $DIST/${VERSION}/${NAME}-${VERSION}.tar.gz $DIST/latest/${NAME}.tar.gz

            tar cvJf $DIST/${VERSION}/${NAME}-${VERSION}.tar.xz .
            cp $DIST/${VERSION}/${NAME}-${VERSION}.tar.xz $DIST/latest/${NAME}.tar.xz
        else
            NAME=$(echo $NAME | sed 's/'${EXT}'//g')
            zip -r $DIST/${VERSION}/${NAME}-${VERSION}.zip *
            cp $DIST/${VERSION}/${NAME}-${VERSION}.zip $DIST/latest/${NAME}.zip
        fi
    )
done