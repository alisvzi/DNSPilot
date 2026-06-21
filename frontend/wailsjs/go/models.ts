export namespace models {
	
	export class DNSInfo {
	    adapter_name: string;
	    dns_servers: string[];
	    is_active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DNSInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.adapter_name = source["adapter_name"];
	        this.dns_servers = source["dns_servers"];
	        this.is_active = source["is_active"];
	    }
	}

}

