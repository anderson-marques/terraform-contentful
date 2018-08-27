FROM hashicorp/terraform:full

RUN go get github.com/hashicorp/terraform/terraform \
    && go get github.com/hashicorp/terraform/helper/resource \
    && go get github.com/hashicorp/terraform/helper/schema \
    && go get github.com/contentful-labs/contentful-go

ENTRYPOINT [ "terraform" ]
