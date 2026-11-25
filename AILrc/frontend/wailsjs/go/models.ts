export namespace main {
	
	export class AppConfig {
	    fontSize: number;
	    fontFamily: string;
	    fontWeight: number;
	    fontColor: string;
	    karaokeColor: string;
	    strokeColor: string;
	    bgOpacity: number;
	    textOpacity: number;
	    karaokeMode: string;
	    lyricEffect: string;
	    effectStrength: number;
	    effectSpeed: number;
	    textAlign: string;
	    letterSpacing: number;
	    glowIntensity: number;
	    borderRadius: number;
	    backdropType: string;
	    showNextLine: boolean;
	    isLocked: boolean;
	    windowWidth: number;
	    windowX: number;
	    windowY: number;
	    windowPositionSaved: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fontSize = source["fontSize"];
	        this.fontFamily = source["fontFamily"];
	        this.fontWeight = source["fontWeight"];
	        this.fontColor = source["fontColor"];
	        this.karaokeColor = source["karaokeColor"];
	        this.strokeColor = source["strokeColor"];
	        this.bgOpacity = source["bgOpacity"];
	        this.textOpacity = source["textOpacity"];
	        this.karaokeMode = source["karaokeMode"];
	        this.lyricEffect = source["lyricEffect"];
	        this.effectStrength = source["effectStrength"];
	        this.effectSpeed = source["effectSpeed"];
	        this.textAlign = source["textAlign"];
	        this.letterSpacing = source["letterSpacing"];
	        this.glowIntensity = source["glowIntensity"];
	        this.borderRadius = source["borderRadius"];
	        this.backdropType = source["backdropType"];
	        this.showNextLine = source["showNextLine"];
	        this.isLocked = source["isLocked"];
	        this.windowWidth = source["windowWidth"];
	        this.windowX = source["windowX"];
	        this.windowY = source["windowY"];
	        this.windowPositionSaved = source["windowPositionSaved"];
	    }
	}
	export class DOMRect {
	    x: number;
	    y: number;
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new DOMRect(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}
	export class LyricWord {
	    time: number;
	    text: string;
	
	    static createFrom(source: any = {}) {
	        return new LyricWord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.text = source["text"];
	    }
	}
	export class LyricLine {
	    time: number;
	    endTime: number;
	    mainText: string;
	    subText: string;
	    words?: LyricWord[];
	
	    static createFrom(source: any = {}) {
	        return new LyricLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.endTime = source["endTime"];
	        this.mainText = source["mainText"];
	        this.subText = source["subText"];
	        this.words = this.convertValues(source["words"], LyricWord);
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
	
	export class MusicInfo {
	    FileName: string;
	    TrackNumber: number;
	
	    static createFrom(source: any = {}) {
	        return new MusicInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FileName = source["FileName"];
	        this.TrackNumber = source["TrackNumber"];
	    }
	}
	export class PlayerState {
	    Position: number;
	    State: number;
	
	    static createFrom(source: any = {}) {
	        return new PlayerState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Position = source["Position"];
	        this.State = source["State"];
	    }
	}

}

