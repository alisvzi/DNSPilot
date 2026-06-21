export namespace models {
	
	export class BenchmarkProfile {
	    id: string;
	    name: string;
	    description: string;
	    domains: string[];
	    attempts: number;
	    timeout_ms: number;
	
	    static createFrom(source: any = {}) {
	        return new BenchmarkProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.domains = source["domains"];
	        this.attempts = source["attempts"];
	        this.timeout_ms = source["timeout_ms"];
	    }
	}
	export class BenchmarkResult {
	    id: string;
	    server_id: string;
	    server_name: string;
	    profile_id: string;
	    profile_name: string;
	    latency_ms: number;
	    jitter_ms: number;
	    success_rate: number;
	    packet_loss: number;
	    score: number;
	    attempts: number;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new BenchmarkResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.server_id = source["server_id"];
	        this.server_name = source["server_name"];
	        this.profile_id = source["profile_id"];
	        this.profile_name = source["profile_name"];
	        this.latency_ms = source["latency_ms"];
	        this.jitter_ms = source["jitter_ms"];
	        this.success_rate = source["success_rate"];
	        this.packet_loss = source["packet_loss"];
	        this.score = source["score"];
	        this.attempts = source["attempts"];
	        this.created_at = source["created_at"];
	    }
	}
	export class DNSServer {
	    id: string;
	    name: string;
	    primary_ipv4: string;
	    secondary_ipv4: string;
	    primary_ipv6: string;
	    secondary_ipv6: string;
	    provider: string;
	    description: string;
	    category: string;
	    tags: string[];
	    is_custom: boolean;
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
	        this.primary_ipv6 = source["primary_ipv6"];
	        this.secondary_ipv6 = source["secondary_ipv6"];
	        this.provider = source["provider"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.tags = source["tags"];
	        this.is_custom = source["is_custom"];
	        this.enabled = source["enabled"];
	    }
	}
	export class CustomDNSList {
	    id: string;
	    name: string;
	    description: string;
	    servers: DNSServer[];
	
	    static createFrom(source: any = {}) {
	        return new CustomDNSList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.servers = this.convertValues(source["servers"], DNSServer);
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
	
	export class NetworkAdapter {
	    id: string;
	    name: string;
	    description: string;
	    mac: string;
	    ipv4: string;
	    gateway: string;
	    dns_servers: string[];
	    is_up: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NetworkAdapter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.mac = source["mac"];
	        this.ipv4 = source["ipv4"];
	        this.gateway = source["gateway"];
	        this.dns_servers = source["dns_servers"];
	        this.is_up = source["is_up"];
	    }
	}
	export class Settings {
	    locale: string;
	    theme: string;
	    auto_apply_fastest: boolean;
	    benchmark_attempts: number;
	    benchmark_timeout_ms: number;
	    test_domain: string;
	    selected_adapter_id: string;
	    last_profile_id: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.locale = source["locale"];
	        this.theme = source["theme"];
	        this.auto_apply_fastest = source["auto_apply_fastest"];
	        this.benchmark_attempts = source["benchmark_attempts"];
	        this.benchmark_timeout_ms = source["benchmark_timeout_ms"];
	        this.test_domain = source["test_domain"];
	        this.selected_adapter_id = source["selected_adapter_id"];
	        this.last_profile_id = source["last_profile_id"];
	    }
	}

}

