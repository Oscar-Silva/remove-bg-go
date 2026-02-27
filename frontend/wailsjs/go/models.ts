export namespace main {
	
	export class ModelConfig {
	    id: string;
	    name: string;
	    sizeMB: number;
	    ram: string;
	    speed: string;
	    quality: string;
	    isDefault: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ModelConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.sizeMB = source["sizeMB"];
	        this.ram = source["ram"];
	        this.speed = source["speed"];
	        this.quality = source["quality"];
	        this.isDefault = source["isDefault"];
	    }
	}

}

