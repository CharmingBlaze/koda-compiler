export namespace api {
	
	export class Diagnostic {
	    path: string;
	    line: number;
	    col: number;
	    message: string;
	    severity: string;
	
	    static createFrom(source: any = {}) {
	        return new Diagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.line = source["line"];
	        this.col = source["col"];
	        this.message = source["message"];
	        this.severity = source["severity"];
	    }
	}

}

export namespace main {
	
	export class DirEntry {
	    name: string;
	    rel: string;
	    isDir: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DirEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.rel = source["rel"];
	        this.isDir = source["isDir"];
	    }
	}

}

