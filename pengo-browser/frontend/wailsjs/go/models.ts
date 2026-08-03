export namespace pengo {
	
	export class Response {
	    Version: string;
	    Status: string;
	    ContentLength: number;
	    ContentType: string;
	    Body: number[];
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Version = source["Version"];
	        this.Status = source["Status"];
	        this.ContentLength = source["ContentLength"];
	        this.ContentType = source["ContentType"];
	        this.Body = source["Body"];
	    }
	}

}

