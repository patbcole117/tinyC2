package ui

import (
	"math/rand"
)

func GetRandomBanner() string {
	return getBanner(rand.Intn(len(banners)))
}
func getBanner(i int) string {
	return banners[i]
}

var banners = [...]string{
	`   **   **                     ******   **** 
	/**  //            **   **  **////** */// *
   ****** ** *******  //** **  **    // /    /*
  ///**/ /**//**///**  //***  /**          *** 
	/**  /** /**  /**   /**   /**         *//  
	/**  /** /**  /**   **    //**    ** *     
	//** /** ***  /**  **      //****** /******
	 //  // ///   //  //        //////  ////// `,
	`::::::::::: ::::::::::: ::::    ::: :::   :::  ::::::::   ::::::::  
	 :+:         :+:     :+:+:   :+: :+:   :+: :+:    :+: :+:    :+: 
	 +:+         +:+     :+:+:+  +:+  +:+ +:+  +:+              +:+  
	 +#+         +#+     +#+ +:+ +#+   +#++:   +#+            +#+    
	 +#+         +#+     +#+  +#+#+#    +#+    +#+          +#+      
	 #+#         #+#     #+#   #+#+#    #+#    #+#    #+#  #+#       
	 ###     ########### ###    ####    ###     ########  ########## `,
	`######    ####    ##  ##   ##  ##    ####     ####   
	 ##       ##     ### ##   ##  ##   ##  ##   ##  ##  
	 ##       ##     ######   ##  ##   ##           ##  
	 ##       ##     ######    ####    ##          ##   
	 ##       ##     ## ###     ##     ##        ##     
	 ##       ##     ##  ##     ##     ##  ##   ##      
	 ##      ####    ##  ##     ##      ####    ######  `,
	`     >=>                               >=>             
	 >=>    >>                      >=>   >=>  >=>>=>  
   >=>>==>     >==>>==>  >=>   >=> >=>        >>   >=> 
	 >=>   >=>  >=>  >=>  >=> >=>  >=>            >=>  
	 >=>   >=>  >=>  >=>    >==>   >=>           >=>   
	 >=>   >=>  >=>  >=>     >=>    >=>   >=>  >=>     
	  >=>  >=> >==>  >=>    >=>       >===>   >======> `,
	`      mm      db                           .g8"""bgd          
	  MM                                 .dP'      M          
	mmMMmm   7MM   7MMpMMMb.   7M'    MF'dM'          pd*"*b. 
	  MM      MM    MM    MM    VA   ,V  MM          (O)   j8 
	  MM      MM    MM    MM     VA ,V   MM.             ,;j9 
	  MM      MM    MM    MM      VVV     Mb.     ,'  ,-='    
	   Mbmo .JMML..JMML  JMML.    ,V        "bmmmd'  Ammmmmmm 
								 ,V                           
	OOb"                            `,
	`             __  __  
	|_ .  _     /     _) 
	|_ | | ) \/ \__  /__ 
			 /           `,
	`                88                              ,ad8888ba,    ad888888b,  
	     ,d     ""                             d8"'     "8b  d8"     "88  
	     88                                   d8'                    a8P  
       MM88MMM  88  8b,dPPYba,   8b       d8  88                  ,d8P"   
	     88     88  88P'    "8a   8b     d8'  88                a8P"      
	     88     88  88       88    8b   d8'   Y8,             a8P'        
	     88,    88  88       88     8b,d8'     Y8a.    .a8P  d8"          
 	    "Y888  88  88       88      Y88'        "Y8888Y"'   88888888888  
 							   d8'                              
d8'`,
}