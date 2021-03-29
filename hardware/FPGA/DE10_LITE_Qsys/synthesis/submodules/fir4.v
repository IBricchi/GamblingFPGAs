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
	output	[31:0]	resultc,

	reg [31:0] result; // = dataa[15:0] << datab[15:0];
	reg [31:0] resultc;
	reg [31:0] all;
	reg [31:0] a;
	always @ (dataa) begin
	 	result <= dataa;
		resultc <= dataa;	
	end

endmodule
