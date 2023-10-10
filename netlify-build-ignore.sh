if  [ $COMMIT_REF == $CACHED_COMMIT_REF ] ; then exit 0 ; fi

case ${SITE_NAME} in 

    nitrodocs)
    git diff --quiet $CACHED_COMMIT_REF $COMMIT_REF -- docs mkdocs.yml
    ;;

    nitro-payment-demo)
    git diff --quiet $CACHED_COMMIT_REF $COMMIT_REF -- ./packages/payment-proxy-client
    ;;

esac