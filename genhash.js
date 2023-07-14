// Copied from phoenixbuilder
const fs=require("fs");
const crypto=require("crypto");

function sha256(data){
	return crypto.createHash("sha256").update(data).digest("hex");
}

let hashes={};
let files=fs.readdirSync("build");
for(let i of files){
	hashes[i]=sha256(fs.readFileSync(`build/${i}`));
}
fs.writeFileSync("build/hashes.json",JSON.stringify(hashes,null,"\t"));