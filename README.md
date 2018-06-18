# tagrr

Taggr is a simple, transactional tags database, designed to make it easy to tag compute resources in the absence of an infrastructure API. Its API is similar to the tags APIs exposed by various public cloud services ([AWS](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html), [GCP](https://cloud.google.com/compute/docs/labeling-resources), [Azure](https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-using-tags), [OpenStack](https://developer.openstack.org/api-ref/compute/#server-tags-servers-tags)).

Use tagrr to get and set tags. Tags are keys with optional values. Keys and values are UTF-8 encoded strings. Many tagrr processes can concurrently read the db. Only one tagrr process can write to the db, and not while any other tagrr processes are reading them. Output is sorted in ascending lexicographical order by key.

## Installation

    $ go get github.com/dlipovetsky/tagrr

## Usage

### Set tags
    $ tagrr set foo=1 bar=2 baz=3

### Get tags

    $ tagrr get baz bar foo
    bar:2
    baz:3
    foo:1

### Get all tags matching a prefix

    $ tagrr get --prefix ba
    bar:2
    baz:3

### Get all tags

    $ tagrr get --all
    bar:2
    baz:3
    foo:1

### Get tags in JSON

    $ tagrr get --all --output json
    {
        "bar": "2",
        "baz": "3",
        "foo": "1",
    }

### Unset (delete) tags

    $ tagrr unset baz foo
