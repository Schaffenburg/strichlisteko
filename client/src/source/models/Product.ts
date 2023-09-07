import axios from "axios";
import Api from "@/source/api";

export default class Product {
    public id: number;
    public name :string;
    public imageId :string = "";
    public stock :number = 0;
    public ean :string;
    private price :number = 0;
    public boxSize :number = 0;
    public amount :string = "none";
    public note :string = "";

    constructor(id:number, name:string, imageId:string = "", stock:number = 0, ean:string = "", price:number = 0, boxSize:number = 0, amount:string = "", note:string = "") {
        this.id=id;
        this.name=name;
        this.imageId=imageId;
        this.stock=stock;
        this.ean=ean;
        this.price=price;
        this.boxSize=boxSize;
        this.amount=amount;
        this.note=note;
    }

    getPrice(){
        return this.price/100;
    }

    /**
     * Returns all available products
     */
    static async all(){
        const api = new Api();
        const resp = await api.getAllProducts();
        if(resp.data){
            const products: Product[] = [];
            for(const prod of resp.data){
                products.push(new Product(prod.id,prod.name,prod.image,prod.stock,prod.EAN,prod.price,prod.box_size,prod.amount,prod.note))
            }
            return products;
        }
        return []
    }

    public async buy(userId:number){
        const api = new Api();
        const resp = await api.buyProduct(userId,this.id);
        if(resp.data){
            if(resp.data.info){
                return true
            }
        }
        return false;
    }



}
