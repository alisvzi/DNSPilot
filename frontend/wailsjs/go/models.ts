export namespace models {
	
	export class DNSServer {
	    id: string;
	    name: string;
	    primary_ipv4: string;
	    secondary_ipv4: string;
	    provider: string;
	    description: string;
	    category: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DNSServer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.primary_ipv4 = source["primary_ipv4"];
	        this.secondary_ipv4 = source["secondary_ipv4"];
	        this.provider = source["provider"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.enabled = source["enabled"];
	    }
	}

}

