deps = {
    "headers": {
        "flatbuffers": {
            "url": "https://github.com/google/flatbuffers/archive/v1.8.0.zip",
            "hash": "ab03443a63bda19954dee6c72ac38d3af5b4579a",
            "headers": {
                "srcs": [":get-flatbuffers"],
                "dir": "flatbuffers-1.8.0/include"
            }
        },
        "rapidxml": {
            "url": "https://github.com/odict/rapidxml/archive/v1.13.zip",
            "hash": "963d7782580d081d9fa5966cdca4a6d89ae0285a",
            "headers": {
                "dir": "rapidxml-1.13",
                "srcs": [":get-rapidxml"]
            }
        }
    },
    "libraries": {
        "boost": {
            "url": "https://dl.bintray.com/boostorg/release/1.66.0/source/boost_1_66_0.zip",
            "hash": "68660effcb94170af7ed2e8871749eee9a031103",
            "build": {
                "cmd": [
                    "cp -r $SRCS/* $OUT",
                    "cd $OUT/boost_1_66_0",
                    "./bootstrap.sh --prefix=./",
                    "./b2 install link=static"
                ]
            },
            "headers": {
                "dir": "boost_1_66_0/include",
                "srcs": [":build-boost"]
            },
            "libraries": [
                {
                    "name": "libboost_filesystem.a",
                    "dir": "boost_1_66_0/lib"
                },
                {
                    "name": "libboost_system.a",
                    "dir": "boost_1_66_0/lib"
                }
            ]
        },
        "snappy": {
            "url": "https://github.com/google/snappy/archive/1.1.7.zip",
            "hash": "43d6626c62c961171b301eb934fd05fc26ae28a8",
            "headers": {
                "srcs": [':get-snappy', ':build-snappy'],
            },
            "libraries": [{
                "name": "libsnappy.a",
            }],
            "build": {
                "cmd": [
                    "cd $SRCS/snappy-1.1.7",
                    "mkdir -p build",
                    "cd build",
                    "cmake ../",
                    "make",
                    "cd ..",
                    "cp -r ./build/* $OUT"
                ]
            }
        },
        "clownfish": {
            "url": "https://github.com/apache/lucy-clownfish/archive/rel/v0.6.2.zip",
            "hash": "c71da58ec77705a465202621f6686bf829ebfde5",
            "headers": {
                "srcs": [':build-clownfish'],
                "dir": "lucy-clownfish-rel-v0.6.2/runtime/c/autogen/include",
                "rawcopy": True
            },
            "libraries": [{
                "dir": "lucy-clownfish-rel-v0.6.2/runtime/c",
                "name": "libclownfish.0.6.2.dylib"
            }],
            "build": {
                "cmd": [
                    "cp -r $SRCS/* $OUT",
                    "cd $OUT/lucy-clownfish-rel-v0.6.2/compiler/c",
                    "./configure",
                    "make",
                    "make test",
                    "cd ../../runtime/c",
                    "./configure",
                    "make",
                    "make test"
                ]
            }
        },
        "lucy": {
            "url": "https://github.com/odict/lucy/archive/rel/v0.6.1.zip",
            "hash": "dc4305666f86ad99d5088e4f931ebf0c35d1a155",
            "headers": {
                "srcs": [':build-lucy'],
                "dir": "get-lucy/lucy-rel-v0.6.1/c/autogen/include",
                "rawcopy": True
            },
            "libraries": [{
                "dir": "get-lucy/lucy-rel-v0.6.1/c",
                "name": "liblucy.dylib"
            }],
            "build": {
                "srcs": [':get-lucy', ':build-clownfish'],
                "cmd": [
                    "cp -r $SRCS/* $OUT",
                    "source $OUT/lucy-clownfish-rel-v0.6.2/devel/bin/setup_env.sh",
                    "cd $OUT/get-lucy/lucy-rel-v0.6.1/c",
                    "./configure",
                    "make",
                    "make test"
                ]
            }
        }
    }
}