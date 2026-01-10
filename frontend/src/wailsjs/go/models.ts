export namespace main {
	
	export class GitBlameLine {
	    line: number;
	    hash: string;
	    author: string;
	    date: string;
	    message: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new GitBlameLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.line = source["line"];
	        this.hash = source["hash"];
	        this.author = source["author"];
	        this.date = source["date"];
	        this.message = source["message"];
	        this.content = source["content"];
	    }
	}
	export class GitFileStatus {
	    filename: string;
	    status: string;
	    staged: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GitFileStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.status = source["status"];
	        this.staged = source["staged"];
	    }
	}
	export class GitRemoteInfo {
	    name: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new GitRemoteInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.url = source["url"];
	    }
	}
	export class GitStash {
	    id: string;
	    message: string;
	    branch: string;
	    date: string;
	
	    static createFrom(source: any = {}) {
	        return new GitStash(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.message = source["message"];
	        this.branch = source["branch"];
	        this.date = source["date"];
	    }
	}

}

