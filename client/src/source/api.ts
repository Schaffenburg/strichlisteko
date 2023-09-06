import config from "../config";
import axios from "axios";
export default class Api{
    static baseUrl = config.baseUrl

    private Get(path:string){
        return axios.get(Api.baseUrl+path)
    }
    private Post(path:string, data:object){
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
}
