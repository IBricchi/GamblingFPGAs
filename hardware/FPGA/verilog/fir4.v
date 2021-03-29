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
	reg [8:0] stage7;
	reg [8:0] stage8;
	reg [8:0] stage9;
	reg [8:0] stage10;
	reg [8:0] stage11;
	reg [8:0] stage12;
	reg [8:0] stage13;
	reg [8:0] stage14;
	reg [8:0] stage15;
	reg [8:0] stage16;
	reg [8:0] stage17;
	reg [8:0] stage18;
	reg [8:0] stage19;
	reg [8:0] stage20;
	always @ (dataa) begin	
		//stage25 <= stage24;		
		
		stage20 <= stage19;
		stage19 <= stage18;		
		stage18 <= stage17;
		stage17 <= stage16;		
		stage16 <= stage15;
		stage15 <= stage14;		
		stage14 <= stage13;
		stage13 <= stage12;		
		stage12 <= stage11;
		stage11 <= stage10;		
		stage10 <= stage9;
		stage9 <= stage8;		
		stage8 <= stage7;
		stage7 <= stage6;		
		stage6 <= stage5;
		stage5 <= stage4;		
		stage4 <= stage3;
		stage3 <= stage2;		
		stage2 <= stage1;		
		stage1 <= dataa;
		
	 	result <= (14*dataa) - (3*stage1)  - (8*stage2) + (15*stage3) - (6*stage4) - (14*stage6)  + (26*stage7) - (5*stage8) - (50*stage9) +(112*stage10) 				- (50*stage11) - (5*stage12) + (26*stage14) - (14*stage15) - (6*stage16) + (15*stage17) - (8*stage18) - (3*stage19) + (14*stage20);

		resultc <= dataa;	
	end

endmodule
