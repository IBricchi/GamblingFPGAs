// synopsys translate_off
'timescale 1 ps / 1 ps
// synopsys translate_on
module bet (
	aclr,
	clk_en,
	clock,
	dataa,
	datab,
	datac,
	result
);

	input	aclr;
	input	clk_en;
	input	clock;
	input	[31:0]	dataa; 	// filtered accelerometer 
	input	[31:0]	datab; 	// Total money
	input 	[ 5:0]	datac; 	// switch 0-5
	output	[31:0]	result;	

	//reg 	[31:0] 	result; // = dataa[15:0] << datab[15:0];
	
	// Scale factors & accelerometer 
	reg 	[31:0]	range = 45;
	reg	[31:0]	min = 5;
	reg	[31:0]	xval;

	// Tracking 
	reg	[ 5:0]	segValue = 0;

	

	// Bet Money
	reg	[3:0]	b_digvalue  [5:0];
	

	//Total monet to play with 
	reg	[3:0]	m_digvalue [5:0]; 

	

	//24'h250913;

	// Intermediates
	reg		Qmax;
	reg	[31:0]	intermediateA;
	reg	[31:0]	intermediateB;
	
	assign result[3:0] = b_digvalue[segValue];
//assign
	//intial m_digvalue[0] =  datab[3:0];		
	
	always @ ( * ) begin
	
		m_digvalue[0] <=  4'd5;	

		case (datac[segValue])
			1'b0	:	begin
						b_digvalue[segValue] <= m_digvalue[segValue] * dataa/60;
						//b_digvalue[segValue] <= (m_digvalue[segValue] - ((((m_digvalue[segValue]*100/range)*dataa) - ((m_digvalue[segValue]*100/range)*min))/100));	


					end

			
			1'b1	:	b_digvalue[segValue] <= b_digvalue[segValue];
			default	:	Qmax <= 1'b0;
				
		endcase
	end
endmodule
