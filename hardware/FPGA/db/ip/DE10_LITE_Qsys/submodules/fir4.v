// synopsys translate_off
'timescale 1 ps / 1 ps
// synopsys translate_on
module fir (
	aclr,
	clk_en,
	clock,
	dataa,
	datab,
	result,
	resultc
);

	input	aclr;
	input	clk_en;
	input	clock;
	input	[8:0]	dataa;
	input	[1:0]	datab;
	output	[31:0]	result;
	output	[31:0]	resultc;

	reg [31:0] result; // = dataa[15:0] << datab[15:0];
	reg [31:0] resultc;
	reg [8:0] stage1;
	reg [8:0] stage2;
	reg [8:0] stage3;

	reg [8:0] stage4;
	reg [8:0] stage5;
	reg [8:0] stage6;

	always @ (dataa) begin	
		stage3 <= stage2;		
		stage2 <= stage1;		
		stage1 <= dataa;
		
	 	result <= (dataa) + (stage1) + (stage2) + (stage3) + (stage4) + (stage5) + (stage6);
		resultc <= dataa;	
	end

endmodule
