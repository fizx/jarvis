namespace go reddit.{{project}}

include "baseplate.thrift"

struct EchoRequest { 
  1: string content;
} 

struct EchoResponse { 
  1: string content;
}

service {{class}}Service extends baseplate.BaseplateService {
    EchoResponse echo(
      1: EchoRequest request,
    ) throws ();
}