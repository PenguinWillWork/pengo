export namespace protocol {
	
	export class PengoResponse {
	    Version: string;
	    Status: string;
	    ContentLength: number;
	    Body: string;
	
	    static createFrom(source: any = {}) {
	        return new PengoResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Version = source["Version"];
	        this.Status = source["Status"];
	        this.ContentLength = source["ContentLength"];
	        this.Body = source["Body"];
	    }
	}

}

