export namespace database {
	
	export class Rules {
	    type: string;
	    ip: string;
	    port: string;
	    username?: string;
	    password?: string;
	
	    static createFrom(source: any = {}) {
	        return new Rules(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.ip = source["ip"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}

}

export namespace main {
	
	export class CurrentRule {
	    type: string;
	    ip: string;
	    port: string;
	    needAuth: boolean;
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new CurrentRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.ip = source["ip"];
	        this.port = source["port"];
	        this.needAuth = source["needAuth"];
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}

}

