import config from "../config";
import axios from "axios";
import Product from "@/source/models/Product";
export default class Api{
    static baseUrl = config.baseUrl

    private Get(path:string){
        return axios.get(Api.baseUrl+path)
    }
    private Post(path:string, data:object|string){
        return axios.post(Api.baseUrl+path,data)
    }

    /**
     * Get all available products
     */
    public getAllProducts(){
        return this.Get('/products');
    }

    public getAllUsers(){
        return this.Get('/users');
    }

    public buyProduct(userId:number,productId:number){
        return this.Post('/user/'+userId,{action:"buy",product:productId});
    }

    public createProduct(product:Product){

    }

    public editProduct(product:Product){

    }

    public saveImage(imageData:string){
        return this.Post('/img/new',imageData)
    }
}
