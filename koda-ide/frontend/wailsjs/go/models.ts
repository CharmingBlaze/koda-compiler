export namespace api {
	
	export class Diagnostic {
	    path: string;
	    line: number;
	    col: number;
	    message: string;
	    hint?: string;
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
	        this.hint = source["hint"];
	        this.severity = source["severity"];
	    }
	}
	export class DocPage {
	    rel: string;
	    title: string;
	    category: string;
	    beginner: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DocPage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rel = source["rel"];
	        this.title = source["title"];
	        this.category = source["category"];
	        this.beginner = source["beginner"];
	    }
	}
	export class SDKLine {
	    ok: boolean;
	    label: string;
	    detail: string;
	    fix: string;
	
	    static createFrom(source: any = {}) {
	        return new SDKLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.label = source["label"];
	        this.detail = source["detail"];
	        this.fix = source["fix"];
	    }
	}
	export class SDKStatus {
	    ok: boolean;
	    version: string;
	    installDir: string;
	    stdlibDir: string;
	    lines: SDKLine[];
	
	    static createFrom(source: any = {}) {
	        return new SDKStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.version = source["version"];
	        this.installDir = source["installDir"];
	        this.stdlibDir = source["stdlibDir"];
	        this.lines = this.convertValues(source["lines"], SDKLine);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
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

