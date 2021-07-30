package main

import (
	"log"
	"net/http"

	//"github.com/pusher/pusher-router-go"
	//"encoding/gob"

	"github.com/gorilla/mux"
	//"github.com/gorilla/securecookie"
	//"github.com/gorilla/sessions"
	"./controllers"
)

func initRouter() {
	router := mux.NewRouter()

	controllers.Init()
	port := "3000"

	//log.Println("Eviroment test: ", os.Getenv("MESDBUSER"))

	log.Println("Server started on: router://localhost:" + port)
	//----------ROUTES--------------------
	/*
		router.Handle("/public/", //final url can be anything
		router.StripPrefix("/public/",
		   router.FileServer(router.Dir("public"))))
	*/
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	router.HandleFunc("/", controllers.Index)
	router.HandleFunc("/login", controllers.Login)
	router.HandleFunc("/logout", controllers.Logout)
	router.HandleFunc("/forbidden", controllers.Forbidden)
	router.HandleFunc("/getMetaUser", controllers.GetMetaUser)
	router.HandleFunc("/getUserTypeAccesor", controllers.GetUserTypeAccesor)

	router.HandleFunc("/reportFailures", controllers.ReportFailures)

	/*---------------------------*/
	router.HandleFunc("/Products", controllers.Products)
	router.HandleFunc("/getProducts", controllers.GetProducts)
	router.HandleFunc("/EditProduct", controllers.EditProduct)
	router.HandleFunc("/NewProduct", controllers.NewProduct)
	router.HandleFunc("/UpdateProduct", controllers.UpdateProduct)
	router.HandleFunc("/InsertProduct", controllers.InsertProduct)
	router.HandleFunc("/DeleteProduct", controllers.DeleteProduct)
	router.HandleFunc("/getProductBy", controllers.GetProductBy)

	/*-----------------------Presentations----------------------------*/
	router.HandleFunc("/Presentations", controllers.Presentations)
	router.HandleFunc("/getPresentations", controllers.GetPresentations)
	router.HandleFunc("/getPresentationBy", controllers.GetPresentationBy)

	router.HandleFunc("/NewPresentation", controllers.NewPresentation)
	router.HandleFunc("/InsertPresentation", controllers.InsertPresentation)
	router.HandleFunc("/editPresentation", controllers.EditPresentation)
	router.HandleFunc("/updatePresentation", controllers.UpdatePresentation)
	router.HandleFunc("/deletePresentation", controllers.DeletePresentation)

	/*-----------------------LTC----------------------------*/
	router.HandleFunc("/lTC", controllers.Ltc)
	router.HandleFunc("/LineTimeClassification", controllers.LineTimeClassification)
	router.HandleFunc("/getLTC", controllers.GetLTC)
	router.HandleFunc("/newLTC", controllers.NewLtc)
	router.HandleFunc("/InsertLTC", controllers.InsertLTC)
	router.HandleFunc("/editLTC", controllers.Editltc)
	router.HandleFunc("/updateltc", controllers.Updateltc)
	router.HandleFunc("/deleteLTC", controllers.Deleteltc)

	router.HandleFunc("/getLTCLineBy", controllers.GetLTCLineBy)
	router.HandleFunc("/LTCLine", controllers.LTCLine)

	router.HandleFunc("/insertLTCLine", controllers.InsertLTCLine)
	router.HandleFunc("/deleteLTCLine", controllers.DeleteLTCLine)

	/*-----------------------Subclassification----------------------------*/
	router.HandleFunc("/SubClassification", controllers.SubClassification)
	router.HandleFunc("/getSub", controllers.GetSub)
	router.HandleFunc("/getSubE", controllers.GetSubE)
	router.HandleFunc("/InsertSub", controllers.InsertSub)
	router.HandleFunc("/newSubClassification", controllers.NewSubClassification)
	router.HandleFunc("/editSubClassification", controllers.EditSubClassification)
	router.HandleFunc("/updateSubClassification", controllers.UpdateSubClassification)
	router.HandleFunc("/deleteSubClassification", controllers.DeleteSubClassification)

	/*-----------------------Branch----------------------------*/
	router.HandleFunc("/getBranch", controllers.GetBranch)
	router.HandleFunc("/InsertBranch", controllers.InsertBranch)
	router.HandleFunc("/Branch", controllers.Branch)
	router.HandleFunc("/newBranch", controllers.NewBranch)
	router.HandleFunc("/editBranch", controllers.EditBranch)
	router.HandleFunc("/updateBranch", controllers.UpdateBranch)
	router.HandleFunc("/deleteBranch", controllers.DeleteBranch)

	/*-----------------------Event----------------------------*/
	router.HandleFunc("/getEvent", controllers.GetEvent)
	router.HandleFunc("/getEventBy", controllers.GetEventBy)
	router.HandleFunc("/getEventFilterV00", controllers.GetEventFilterV00)
	router.HandleFunc("/getEventFilterV01", controllers.GetEventFilterV01)
	router.HandleFunc("/InsertEvent", controllers.InsertEvent)
	router.HandleFunc("/Event", controllers.Event)
	router.HandleFunc("/newEvent", controllers.NewEvent)
	router.HandleFunc("/editEvent", controllers.EditEvent)
	router.HandleFunc("/updateEvent", controllers.UpdateEvent)
	router.HandleFunc("/deleteEvent", controllers.DeleteEvent)

	router.HandleFunc("/getExLby", controllers.GetExLby)
	router.HandleFunc("/deleteExL", controllers.DeleteExL)
	router.HandleFunc("/editExL", controllers.EditExL)
	router.HandleFunc("/updateExL", controllers.UpdateExL)
	router.HandleFunc("/updateExLPartial", controllers.UpdateExLPartial)
	/*-------------------Report-----WCM-----------------------------*/
	router.HandleFunc("/selectRegisterWCM", controllers.SelectRegisterWCM)
	router.HandleFunc("/reportEvent", controllers.ReportEvent)
	router.HandleFunc("/InsertReportEvent", controllers.InsertReportEvent)
	router.HandleFunc("/InsertReportEventValidation", controllers.InsertReportEventValidation)

	router.HandleFunc("/reportPlanning", controllers.ReportPlanning)
	router.HandleFunc("/reportPlanningbyWeek", controllers.ReportPlanningbyWeek)
	router.HandleFunc("/reportProduced", controllers.ReportProduced)
	router.HandleFunc("/updateProducedBox", controllers.UpdateProducedBox)
	router.HandleFunc("/updateProducedBoxLog", controllers.UpdateProducedBoxLog)

	router.HandleFunc("/getActualPlanningBy", controllers.GetActualPlanningBy)

	router.HandleFunc("/InsertReportPlanning", controllers.InsertReportPlanning)
	router.HandleFunc("/getHistoricPlanning", controllers.GetHistoricPlanning)
	router.HandleFunc("/insertReportPlanningByWeek", controllers.InsertReportPlanningByWeek)

	/*-------------------Analisis----------------------------------*/
	router.HandleFunc("/consolidado", controllers.Consolidado)
	router.HandleFunc("/HistoricPlanning", controllers.HistoricPlanning)

	/*----------------------Sectors-------------------------------*/
	router.HandleFunc("/Factory", controllers.Factory)
	router.HandleFunc("/newFactory", controllers.NewFactory)
	router.HandleFunc("/insertFactory", controllers.InsertFactory)
	router.HandleFunc("/editFactory", controllers.EditFactory)
	router.HandleFunc("/updateFactory", controllers.UpdateFactory)
	router.HandleFunc("/deleteFactory", controllers.DeleteFactory)

	router.HandleFunc("/Area", controllers.Area)
	router.HandleFunc("/newArea", controllers.NewArea)
	router.HandleFunc("/insertArea", controllers.InsertArea)
	router.HandleFunc("/editArea", controllers.EditArea)
	router.HandleFunc("/updateArea", controllers.UpdateArea)
	router.HandleFunc("/getArea", controllers.GetArea)

	router.HandleFunc("/getAreaBy", controllers.GetAreaBy)
	router.HandleFunc("/deleteArea", controllers.DeleteArea)

	router.HandleFunc("/Line", controllers.Line)
	router.HandleFunc("/getLineE", controllers.GetLineE)
	router.HandleFunc("/newLine", controllers.NewLine)
	router.HandleFunc("/insertLine", controllers.InsertLine)
	router.HandleFunc("/editLine", controllers.EditLine)
	router.HandleFunc("/deleteLine", controllers.DeleteLine)
	router.HandleFunc("/updateLine", controllers.UpdateLine)

	router.HandleFunc("/getLineBy", controllers.GetLineBy)

	/*---------------------OEE--------------------------------*/
	router.HandleFunc("/getOEE", controllers.GetOEE)
	router.HandleFunc("/getRealtiveOEE", controllers.GetRelativeOEEDayData)
	router.HandleFunc("/filterOEE", controllers.FilterOEE)
	router.HandleFunc("/oeeDemo00", controllers.OEEDemo00)
	router.HandleFunc("/oeeDemo01", controllers.OEEDemo01)
	router.HandleFunc("/oeeDemo02", controllers.OEEDemo02)
	router.HandleFunc("/getPlanningV00", controllers.GetPlanningV00)
	router.HandleFunc("/getPlanningV02", controllers.GetPlanningV02)
	router.HandleFunc("/monitor00", controllers.Monitor00)
	router.HandleFunc("/monitorSetup", controllers.MonitorSetup)
	router.HandleFunc("/consolidateOEE", controllers.ConsolidateOEE)
	router.HandleFunc("/getOEEWeek", controllers.GetOEEWeek)
	router.HandleFunc("/getOEEWeekxLineAPI", controllers.GetOEEWeekxLineAPI)
	router.HandleFunc("/getOEEProjectionbyRange", controllers.GetOEEProjectionbyRange)

	router.HandleFunc("/OEEYTDProjection", controllers.OEEYTDProjection)

	router.HandleFunc("/getOEEWeekLineTotal", controllers.GetOEEWeekLineTotal)
	router.HandleFunc("/consolidateOEEWeek", controllers.ConsolidateOEEWeek)

	router.HandleFunc("/ws", controllers.WsEndpoint)
	router.HandleFunc("/wsEvent", controllers.WsEndpointEvent)
	router.HandleFunc("/wsWeight", controllers.WsEndpointWeight)
	/*---------------------AMIS--------------------------------*/
	router.HandleFunc("/consolidatedAMIS", controllers.ConsolidatedAMIS)
	router.HandleFunc("/getAMISV00", controllers.GetAMISV00)
	router.HandleFunc("/calcORbyTurn", controllers.CalcORbyTurn)

	/*---------------------Planning--------------------------------*/
	router.HandleFunc("/planningValidation", controllers.PlanningValidation)
	router.HandleFunc("/editPlanning", controllers.EditPlanning)
	router.HandleFunc("/updatePlanning", controllers.UpdatePlanning)
	router.HandleFunc("/deletePlanning", controllers.DeletePlanning)
	//router.HandleFunc("/DeleteExL", DeleteExL)

	/*-----------------Validation------------------------------------*/
	router.HandleFunc("/validationOEE", controllers.ValidationOEE)
	router.HandleFunc("/eventsTemplate", controllers.EventsTemplate)

	router.HandleFunc("/getDTime", controllers.GetDTime)

	/*---------------------DMS--------------------------*/
	router.HandleFunc("/selectRegister", controllers.SelectRegister)
	router.HandleFunc("/fillSubHeaderPChemical", controllers.FillSubHeaderPChemical)
	router.HandleFunc("/insertSubPChemical", controllers.InsertSubPChemical)
	router.HandleFunc("/fillPChemicalGeneral", controllers.FillPChemicalGeneral)

	router.HandleFunc("/getHeaderBy", controllers.GetHeaderBy)
	router.HandleFunc("/getSensorial_analysis_scale", controllers.GetSensorial_analysis_scale)
	router.HandleFunc("/getDesition_catalog", controllers.GetDesition_catalog)
	router.HandleFunc("/insertPChemicalGeneral", controllers.InsertPChemicalGeneral)
	router.HandleFunc("/selectDMSValidation", controllers.SelectDMSValidation)

	router.HandleFunc("/validationPChemical", controllers.ValidationPChemical)
	router.HandleFunc("/getPChemicalGeneral", controllers.GetPChemicalGeneral)
	router.HandleFunc("/getPChemicalV00", controllers.GetPChemicalV00)
	router.HandleFunc("/consolidatedPChemical", controllers.ConsolidatedPChemical)

	/*---------------------CALREG-022--------------------------*/
	router.HandleFunc("/fillSalsitasControl", controllers.FillSalsitasControl)

	/**------------------------------------*/
	router.HandleFunc("/testPDF", controllers.WeightSigner2)
	/*-----------------------------------------------------*/

	/*------------------------Gestion-----------------------------*/
	router.HandleFunc("/gestion_personal", controllers.Gestion_personal)
	router.HandleFunc("/newUser", controllers.NewUser)
	router.HandleFunc("/getFactory", controllers.GetFactory)
	router.HandleFunc("/getPrivilege", controllers.GetPrivilege)
	router.HandleFunc("/insertUser", controllers.InsertUser)
	router.HandleFunc("/getUsers", controllers.GetUsers)
	router.HandleFunc("/getLineManagers", controllers.GetLineManagers)
	router.HandleFunc("/getMecanics", controllers.GetMecanics)
	router.HandleFunc("/editUser", controllers.EditUser)
	router.HandleFunc("/updateUser", controllers.UpdateUser)
	router.HandleFunc("/deleteUser", controllers.DeleteUser)

	//router.HandleFunc("/fill")
	/*-----------------------------------------------------*/
	router.HandleFunc("/insertSalsitasControl", controllers.InsertSalsitasControl)
	router.HandleFunc("/insertSalsitasWeight_control", controllers.InsertSalsitasWeight_control)
	router.HandleFunc("/consolidatedWeightSalsitas", controllers.ConsolidatedWeightSalsitas)
	router.HandleFunc("/weightSubControlStep", controllers.WeightSubControlStep)
	/*-----------------------------------------------------*/
	router.HandleFunc("/wsS4", controllers.WsEndpointS4)
	router.HandleFunc("/wsS8", controllers.WsEndpointS8)
	router.HandleFunc("/getWeightSalsitas", controllers.GetWeightSalsitas)
	router.HandleFunc("/getWeightAll", controllers.GetWeightAll)

	router.HandleFunc("/weightVerification", controllers.WeightVerification)
	router.HandleFunc("/weightSigner", controllers.WeightSigner2)
	router.HandleFunc("/signProceedVerification", controllers.SignProceedVerification)

	//log.Fatal(router.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil))

	/*-------------------Manual----------------------------------*/
	router.HandleFunc("/manual", controllers.Manual)

	/*-------------------Configurar-Trabajo----------------------------------*/
	router.HandleFunc("/setupJob", controllers.SetupJob)
	router.HandleFunc("/setJob", controllers.SetJob)

	/*----------------Material------------------------*/
	router.HandleFunc("/material", controllers.Material)
	router.HandleFunc("/newMaterial", controllers.NewMaterial)
	router.HandleFunc("/insertMaterial", controllers.InsertMaterial)
	router.HandleFunc("/editMaterial", controllers.EditMaterial)
	router.HandleFunc("/updateMaterial", controllers.UpdateMaterial)
	router.HandleFunc("/deleteMaterial", controllers.DeleteMaterial)
	router.HandleFunc("/getMaterialBy", controllers.GetMaterialBy)
	router.HandleFunc("/updateEquipment", controllers.UpdateEquipment)
	router.HandleFunc("/deleteEquipment", controllers.DeleteEquipment)

	/*-----------------Call-Off---------------------------------*/
	router.HandleFunc("/selectRegisterBodega", controllers.SelectRegisterBodega)
	router.HandleFunc("/requestCalloff", controllers.RequestCalloff)
	router.HandleFunc("/insertRequestCalloff", controllers.InsertRequestCalloff)
	router.HandleFunc("/insertRequestCalloff", controllers.InsertRequestCalloff)
	router.HandleFunc("/updateCalloff", controllers.UpdateCalloff)

	router.HandleFunc("/gestionBodega", controllers.GestionCalloff)
	router.HandleFunc("/gestionCalloff", controllers.GestionCalloff)

	router.HandleFunc("/getCalloffV00", controllers.GetCalloffV00)
	router.HandleFunc("/getCalloffV01", controllers.GetCalloffV01)
	router.HandleFunc("/getCalloffV02", controllers.GetCalloffV02)
	router.HandleFunc("/getCalloffEV00", controllers.GetCalloffEV00)
	router.HandleFunc("/consolidatedCalloff", controllers.ConsolidatedCalloff)
	router.HandleFunc("/gestionSolicitudCalloff", controllers.GestionSolicitudCalloff)

	router.HandleFunc("/wsCalloff", controllers.WsEndpointCalloff)

	/*----------------Equipment------------------------*/
	router.HandleFunc("/getEquipment", controllers.GetEquipment)
	router.HandleFunc("/getEquipmentBy", controllers.GetEquipmentBy)
	router.HandleFunc("/equipment", controllers.Equipment)
	router.HandleFunc("/newEquipment", controllers.NewEquipment)
	router.HandleFunc("/insertEquipment", controllers.InsertEquipment)
	router.HandleFunc("/getEquipment", controllers.GetEquipmentBy)

	/*----------------Priority------------------------*/
	router.HandleFunc("/getPriorityBy", controllers.GetPriorityBy)

	/*--------------------Boletas-----------------------------------------*/
	router.HandleFunc("/selectRegisterTag", controllers.SelectRegisterTag)

	router.HandleFunc("/blueTag", controllers.BlueTag)
	router.HandleFunc("/insertRedTag", controllers.InsertRedTag)

	router.HandleFunc("/redTag", controllers.RedTag)
	router.HandleFunc("/insertBlueTag", controllers.InsertBlueTag)

	router.HandleFunc("/greenTag", controllers.GreenTag)

	router.HandleFunc("/orangeTag", controllers.OrangeTag)
	router.HandleFunc("/insertOrangeTag", controllers.InsertOrangeTag)

	router.HandleFunc("/gestionSolicitudBoletas", controllers.GestionSolicitudBoletas)
	router.HandleFunc("/gestionCurrentTags", controllers.GestionCurrentTags)

	router.HandleFunc("/updateTag", controllers.UpdateTag)
	router.HandleFunc("/setPriotityTag", controllers.SetPriotityTag)
	router.HandleFunc("/getTagsV00", controllers.GetTagsV00)
	router.HandleFunc("/getTagsV02", controllers.GetTagsV02)
	router.HandleFunc("/getQaTagsV00", controllers.GetQaTagsV00)
	router.HandleFunc("/getTagsV01", controllers.GetTagsV01)
	router.HandleFunc("/getTagEV00", controllers.GetTagEV00)

	router.HandleFunc("/closeTag", controllers.CloseTag)
	router.HandleFunc("/consolidatedTags", controllers.ConsolidatedTags)
	router.HandleFunc("/consolidatedQaTags", controllers.ConsolidatedQaTags)

	router.HandleFunc("/wsTag", controllers.WsEndpointTag) //blue red in real time

	router.HandleFunc("/getQa_anomaly_catalog", controllers.GetQa_anomaly_catalog)
	router.HandleFunc("/countOpenTagsBy", controllers.CountOpenTagsBy)

	/*-----------------------GreenTag----------------------------------*/
	router.HandleFunc("/getClass_of_event", controllers.GetClass_of_event)
	router.HandleFunc("/getEvent_cause", controllers.GetEvent_cause)
	router.HandleFunc("/getFrequency_catalog", controllers.GetFrequency_catalog)
	router.HandleFunc("/getSeverity_catalog", controllers.GetSeverity_catalog)
	router.HandleFunc("/getSHE_standard_catalog", controllers.GetSHE_standard_catalog)

	router.HandleFunc("/insertGreenTag", controllers.InsertGreenTag)

	router.HandleFunc("/getSHETagsV00", controllers.GetSHETagsV00)

	router.HandleFunc("/consolidatedSHETags", controllers.ConsolidatedSHETags)

	/*-----------------------Proecess-Temperature-Control--------------------*/
	router.HandleFunc("/processTemperatureControlStep", controllers.ProcessTemperatureControlStep)
	router.HandleFunc("/insertTemperatureControl", controllers.InsertTemperatureControl)
	router.HandleFunc("/consolidatedTemperatureControl", controllers.ConsolidatedTemperatureControl)
	router.HandleFunc("/getTemperatureControlV00", controllers.GetTemperatureControlV00)

	/*----------------------Jaw_teflon_ultrasonic_state-------------------------*/
	router.HandleFunc("/jaw_teflon_ultrasonic_StateStep", controllers.Jaw_teflon_ultrasonic_StateStep)
	router.HandleFunc("/insertJawControl", controllers.InsertJawControl)
	router.HandleFunc("/getJawControlV00", controllers.GetJawControlV00)
	router.HandleFunc("/consolidatedJawControl", controllers.ConsolidatedJawControl)

	/*----------------------Seal_verification_pneumaticpress-------------------------*/
	router.HandleFunc("/fillSealStep", controllers.FillSealStep)
	router.HandleFunc("/insertSealControl", controllers.InsertSealControl)
	router.HandleFunc("/getSealControlV00", controllers.GetSealControlV00)
	router.HandleFunc("/consolidatedSealControl", controllers.ConsolidatedSealControl)

	/*---------------------CRQS y verificacion de alergenos - 192-------------------------*/
	router.HandleFunc("/getCQRS_Category", controllers.GetCQRS_Category)
	router.HandleFunc("/getCRQS_SubCategoryBy", controllers.GetCRQS_SubCategoryBy)
	router.HandleFunc("/CRQSStep", controllers.CRQSStep)

	router.HandleFunc("/insertCRQS", controllers.InsertCRQS)

	router.HandleFunc("/GetCRQSV00", controllers.GetCRQSV00)
	router.HandleFunc("/consolidatedCRQS", controllers.ConsolidatedCRQS)

	/*---------------------CRQS y verificacion de alergenos - 192-------------------------*/
	router.HandleFunc("/getReason_change", controllers.GetReason_change)
	router.HandleFunc("/allergenVerificationStep", controllers.AllergenVerificationStep)
	router.HandleFunc("/insertAllergenVerification", controllers.InsertAllergenVerification)

	router.HandleFunc("/getAllergenVerificationV00", controllers.GetAllergenVerificationV00)
	router.HandleFunc("/consolidatedAllergenVerification", controllers.ConsolidatedAllergenVerification)

	/*---------------------LossTree-------------------------*/
	router.HandleFunc("/lossTree", controllers.LossTree)
	router.HandleFunc("/getLossTreeData", controllers.GetLossTreeData)
	router.HandleFunc("/getLossTreeDataGrid", controllers.GetLossTreeDataGrid)
	router.HandleFunc("/getLossTreeDataGridBy", controllers.GetLossTreeDataGridBy)
	router.HandleFunc("/getByLineLossTreeDataGridBy", controllers.GetByLineLossTreeDataGridBy)

	/*---------------------Batch-------------------------*/
	router.HandleFunc("/reportBatch", controllers.ReportBatch)
	router.HandleFunc("/insertBatchChange", controllers.InsertBatch)
	router.HandleFunc("/getBatchBy", controllers.GetBatchBy)
	router.HandleFunc("/consolidatedBatch", controllers.ConsolidatedBatch)

	/*------------------Job_catalog--------------------*/
	router.HandleFunc("/job_catalog", controllers.Job_catalog)
	router.HandleFunc("/newJob_catalog", controllers.NewJob_catalog)
	router.HandleFunc("/editJob_catalog", controllers.EditJob_catalog)
	router.HandleFunc("/insertJob_catalog", controllers.InsertJob_catalog)
	router.HandleFunc("/updateJob_catalog", controllers.UpdateJob_catalog)
	router.HandleFunc("/deleteJob_catalog", controllers.DeleteJob_catalog)
	router.HandleFunc("/getJob_catalog", controllers.GetJob_catalog)
	router.HandleFunc("/getJob_catalogE", controllers.GetJob_catalogE)

	/*------------------Work-Order--------------------*/
	router.HandleFunc("/work_order", controllers.Work_order)
	router.HandleFunc("/insertWork_order", controllers.InsertWork_order)
	router.HandleFunc("/getMetaWork_order", controllers.GetMetaWork_order)
	router.HandleFunc("/getMetaWork_orderE", controllers.GetMetaWork_orderE)
	router.HandleFunc("/executeWorkOrder", controllers.ExecuteWorkOrder)
	router.HandleFunc("/mantenimiento", controllers.Mantenimiento)
	router.HandleFunc("/setWorkOrderPhase", controllers.SetWorkOrderPhase)
	router.HandleFunc("/setWorkOrderNote", controllers.SetWorkOrderNote)
	router.HandleFunc("/setWorkOrderPhotoBefore", controllers.SetWorkOrderPhotoBefore)
	router.HandleFunc("/setWorkOrderPhotoAfter", controllers.SetWorkOrderPhotoAfter)
	router.HandleFunc("/work_orderLog", controllers.Work_orderLog)
	router.HandleFunc("/getMetaLogWork_order", controllers.GetMetaLogWork_order)
	router.HandleFunc("/getPlannedWorkSaturarion", controllers.GetPlannedWorkSaturarion)
	router.HandleFunc("/getActualWorkSaturarion", controllers.GetActualWorkSaturarion)
	router.HandleFunc("/workOrderSaturation", controllers.WorkOrderSaturation)
	router.HandleFunc("/deleteWorkOrder", controllers.DeleteWorkOrder)
	router.HandleFunc("/relativeWorkOrder", controllers.RelativeWork_order)
	router.HandleFunc("/getMetaWork_orderList", controllers.GetMetaWork_orderList)

	/*------------------AMCalendar--------------------*/
	router.HandleFunc("/Machine", controllers.Machine)
	router.HandleFunc("/getMachineBy", controllers.GetMachineBy)
	router.HandleFunc("/newMachine", controllers.NewMachine)
	router.HandleFunc("/insertMachine", controllers.InsertMachine)
	router.HandleFunc("/getMachine", controllers.GetMachineBy)
	router.HandleFunc("/updateMachine", controllers.UpdateMachine)
	router.HandleFunc("/deleteMachine", controllers.DeleteMachine)
	router.HandleFunc("/getMachineCatalog", controllers.GetMachineCatalog)

	router.HandleFunc("/Component", controllers.Component)
	router.HandleFunc("/newComponent", controllers.NewComponent)
	router.HandleFunc("/insertComponent", controllers.InsertComponent)
	router.HandleFunc("/getComponent", controllers.GetComponent)
	router.HandleFunc("/updateComponent", controllers.UpdateComponent)
	router.HandleFunc("/deleteComponent", controllers.DeleteComponent)
	router.HandleFunc("/getComponentByMachine", controllers.GetComponentByMachine)

	router.HandleFunc("/EPP", controllers.EPP)
	router.HandleFunc("/newEPP", controllers.NewEPP)
	router.HandleFunc("/insertEPP", controllers.InsertEPP)
	router.HandleFunc("/getEPPBy", controllers.GetEPPBy)
	router.HandleFunc("/updateEPP", controllers.UpdateEPP)
	router.HandleFunc("/deleteEPP", controllers.DeleteEPP)

	router.HandleFunc("/AM_Job", controllers.AM_Job)
	router.HandleFunc("/newAM_Job", controllers.NewAM_Job)
	router.HandleFunc("/insertAM_Job", controllers.InsertAM_Job)
	router.HandleFunc("/getAM_JobBy", controllers.GetAM_JobBy)
	router.HandleFunc("/updateAM_Job", controllers.UpdateAM_Job)
	router.HandleFunc("/deleteAM_Job", controllers.DeleteAM_Job)

	router.HandleFunc("/getLILACatalog", controllers.GetLILACatalog)
	router.HandleFunc("/getLILABy", controllers.GetLILABy)

	router.HandleFunc("/AM_instance", controllers.AM_instance)
	router.HandleFunc("/insertAM_instance", controllers.InsertAM_instance)
	router.HandleFunc("/getMetaAM_instance", controllers.GetMetaAM_instance)
	router.HandleFunc("/getMetaLogAM_instance", controllers.GetMetaLogAM_instance)
	router.HandleFunc("/AM_instanceLog", controllers.AM_instanceLog)
	router.HandleFunc("/getMetaLogAM_instanceE", controllers.GetMetaLogAM_instanceE)
	router.HandleFunc("/deleteAM_Intance", controllers.DeleteAM_Intance)
	router.HandleFunc("/ExecuteAM_instance", controllers.ExecuteAM_instance)
	router.HandleFunc("/setAMPhase", controllers.SetAMPhase)
	router.HandleFunc("/setAMNote", controllers.SetAMNote)
	router.HandleFunc("/SetAMFlashPhase", controllers.SetAMFlashPhase)
	router.HandleFunc("/getAM_Stats", controllers.GetAM_Stats)
	router.HandleFunc("/AMStats", controllers.AMStats)

	// ----------------------------------------------------
	router.HandleFunc("/getMetaLogAM_instanceEList", controllers.GetMetaLogAM_instanceEList)

	// --------------------DFOS--------------------------------
	router.HandleFunc("/dfos", controllers.DFOS)

	// --------------------PBI-Reports--------------------------------

	router.HandleFunc("/pbiOEEReport", controllers.PBIOEEReport)
	router.HandleFunc("/pbiEWOReport", controllers.PBIEWOReport)

	// --------------------General Dashboard--------------------------------
	router.HandleFunc("/specificDashBoard", controllers.SpecificDashBoard)
	router.HandleFunc("/getQACurrentMonitor", controllers.GetQACurrentMonitor)
	router.HandleFunc("/getQARelativeMonitor", controllers.GetQARelativeMonitor)

	router.HandleFunc("/getQARelativeMonitorWeekUtilitation", controllers.GetQARelativeMonitorWeekUtilitation)

	router.HandleFunc("/generalDashboard", controllers.GeneralDashboard)
	router.HandleFunc("/dashBoardCounterBy", controllers.DashBoardCounterBy)

	// --------------------CleanDisinfection--------------------------------
	router.HandleFunc("/cleanDisinfection", controllers.CleanDisinfection)
	router.HandleFunc("/getCleanFilter", controllers.GetCleanFilter)
	router.HandleFunc("/getWashingStage", controllers.GetWashingStage)
	router.HandleFunc("/insertCleanDisinfection", controllers.InsertCleanDisinfection)
	router.HandleFunc("/getCleanDisinfectionBy", controllers.GetCleanDisinfectionBy)
	router.HandleFunc("/consolidatedCleanDisinfection", controllers.ConsolidatedCleanDisinfection)
	router.HandleFunc("/getCleanDisinfectionE", controllers.GetCleanDisinfectionE)
	router.HandleFunc("/setCleanState", controllers.SetCleanState)

	// --------------------SalsitasOUtputControl--------------------------------
	router.HandleFunc("/salsitasOutputControl", controllers.SalsitasOutputControl)
	router.HandleFunc("/consolidatedSalsitasOutputControl", controllers.ConsolidatedSalsitasOutputControl)

	router.HandleFunc("/getClean_corrective_action", controllers.GetClean_corrective_action)

	router.HandleFunc("/insertSalsitasOutputControl", controllers.InsertSalsitasOutputControl)
	router.HandleFunc("/getSalsitasOutputControlBy", controllers.GetSalsitasOutputControlBy)

	// --------------------Coding Verification--------------------------------
	router.HandleFunc("/codingVerification", controllers.CodingVerification)
	router.HandleFunc("/insertCodingVerification", controllers.InsertCodingVerification)
	router.HandleFunc("/getCodingVerificationBy", controllers.GetCodingVerificationBy)
	router.HandleFunc("/consolidatedCodingVerification", controllers.ConsolidatedCodingVerification)
	router.HandleFunc("/setCodingVerification", controllers.SetCodingVerification)
	router.HandleFunc("/getCodingVerificationE", controllers.GetCodingVerificationE)

	router.HandleFunc("/getOEEMetaProjectionbyRange", controllers.GetOEEMetaProjectionbyRange)

	// --------------------PackingControl--------------------------------
	router.HandleFunc("/packingControl", controllers.PackingControl)
	router.HandleFunc("/insertPackingControl", controllers.InsertPackingControl)

	router.HandleFunc("/getPackingControlBy", controllers.GetPackingControlBy)
	router.HandleFunc("/consolidatedPackingControl", controllers.ConsolidatedPackingControl)

	// --------------------TemporalStorageProduct--------------------------------

	router.HandleFunc("/temporalStorageProduct", controllers.TemporalStorageProduct)
	router.HandleFunc("/insertTemporalStorageProduct", controllers.InsertTemporalStorageProduct)

	router.HandleFunc("/getTemporalStorageProduct", controllers.GetTemporalStorageProductBy)
	router.HandleFunc("/consolidatedTemporalStorageProduct", controllers.ConsolidatedTemporalStorageProduct)

	// --------------------TemporalStorageProduct--------------------------------

	router.HandleFunc("/TomateCubeBins", controllers.TomateCubeBins)
	router.HandleFunc("/insertTomateCubeBins", controllers.InsertTomateCubeBins)

	router.HandleFunc("/getTomateCubeBinsBy", controllers.GetTomateCubeBinsBy)
	router.HandleFunc("/consolidatedTomateCubeBins", controllers.ConsolidatedTomateCubeBins)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func main() {

	initRouter()

}
