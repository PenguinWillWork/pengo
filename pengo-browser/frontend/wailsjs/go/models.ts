export namespace pengo {
	
	export class Response {
	    Version: string;
	    Status: string;
	    ContentLength: number;
	    Body: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
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

