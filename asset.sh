cd ./resource/static/ && \
    go-assets-builder . -o ../../data/staticfs.go -p data -v StaticFS && \
    cd ../template && \
    go-assets-builder . -o ../../data/templatefs.go -p data -v TemplateFS